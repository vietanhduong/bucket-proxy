package gcs

import "google.golang.org/api/option"

type Option func(*Client)

func WithCredentialsFile(path string) Option {
	return func(c *Client) {
		if path != "" {
			c.opts = append(c.opts, option.WithCredentialsFile(path))
		}
	}
}
