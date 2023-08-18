package pgopher

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type profileCollector struct {
	ctx        context.Context
	logger     slog.Logger
	target     ProfilingTarget
	sink       Sink
	s3Client   *s3.Client
	kubeClient *kubernetes.Clientset
}

func (p profileCollector) Run() {
	p.logger.Info("collecting profile")

	url := fmt.Sprintf("%s?seconds=%d", p.target.URL, int(p.target.Duration.Seconds()))

	req, err := http.NewRequestWithContext(p.ctx, http.MethodGet, url, nil)
	if err != nil {
		p.logger.Error("failed to create http request", slog.String("err", err.Error()))
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		p.logger.Error("failed to collect profile", slog.String("err", err.Error()))
		return
	}

	defer resp.Body.Close()

	buf := &bytes.Buffer{}

	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		p.logger.Error("failed to read body", slog.String("err", err.Error()))
		return
	}

	switch p.sink.Type {
	case "file":
		filePath := filepath.Join(p.sink.FileSinkOptions.Folder, fmt.Sprintf("%s.pgo", p.target.Name))

		file, err := os.Create(filePath)
		if err != nil {
			p.logger.Error("failed to create file for sink", slog.String("err", err.Error()))
			return
		}

		defer file.Close()

		_, err = file.Write(buf.Bytes())
		if err != nil {
			p.logger.Error("failed to write to file sink", slog.String("err", err.Error()), slog.String("file", file.Name()))
			return
		}
	case "s3":
		_, err = p.s3Client.PutObject(p.ctx, &s3.PutObjectInput{
			Bucket: aws.String(p.sink.S3SinkOptions.Bucket),
			Key:    aws.String(fmt.Sprintf("profile=%s/%s.pgo", p.target.Name, p.target.Name)),
			Body:   buf,
		})
		if err != nil {
			p.logger.Error("failed to write to s3 sink", slog.String("err", err.Error()), slog.String("bucket", p.sink.S3SinkOptions.Bucket))
			return
		}
	case "kubernetes":
		name := fmt.Sprintf("pgopher-profile-%s", p.target.Name)

		var profileData []byte

		base64.StdEncoding.Encode(profileData, buf.Bytes())

		client := p.kubeClient.CoreV1().Secrets(p.sink.KubernetesSinkOptions.Namespace)

		secret := core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      name,
				Namespace: p.sink.KubernetesSinkOptions.Namespace,
			},
			StringData: make(map[string]string),
		}

		secret.StringData["profile"] = string(profileData)

		_, err := client.Get(p.ctx, name, meta.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				_, err := client.Create(p.ctx, &secret, meta.CreateOptions{})
				if err != nil {
					slog.Error("failed to create secret", slog.String("err", err.Error()))
					return
				}

				return
			} else {
				slog.Error("failed to get secret", slog.String("err", err.Error()))
				return
			}
		}

		_, err = client.Update(p.ctx, &secret, meta.UpdateOptions{})
		if err != nil {
			slog.Error("failed to update secret", slog.String("err", err.Error()))
			return
		}
	default:
		p.logger.Error("unknown sink", slog.String("sink", p.sink.Type))
		return
	}
}
