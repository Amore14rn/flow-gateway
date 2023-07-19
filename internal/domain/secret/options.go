package secret

import "github.com/google/uuid"

type Option func(*Secret) error
type Options []Option

func NewOptions(opts ...Option) Options {
	return opts
}

func (o *Options) Append(opt Option) {
	*o = append(*o, opt)
}

func WithID(id string) Option {
	return func(u *Secret) error {
		uid, err := uuid.Parse(id)
		if err != nil {
			return err
		}
		u.UUID = uid
		u.SetChanged("UUID", uid)
		return nil
	}
}

func WithSecret(secret string) Option {
	return func(u *Secret) error {
		u.Secret = secret
		u.SetChanged("Secret", secret)
		return nil
	}
}
