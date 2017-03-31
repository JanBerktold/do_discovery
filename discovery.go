package do

import (
	"github.com/digitalocean/godo"
	"context"
	"golang.org/x/oauth2"
)

type DiscoveryClient struct {
	digitalClient *godo.Client
}

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func NewClient(AccessToken string) *DiscoveryClient {
	oauthClient := oauth2.NewClient(oauth2.NoContext, &TokenSource{
		AccessToken: AccessToken,
	})

	return &DiscoveryClient{
		digitalClient: godo.NewClient(oauthClient),
	}
}

type FilterOptions struct {}

func (c *DiscoveryClient) ByTag(tag string, opts *FilterOptions) ([]godo.Droplet, error) {
	// TODO: inefficent, fix this
	result := make([]godo.Droplet, 256)
	index := 0

	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := c.digitalClient.Droplets.ListByTag(context.Background(), tag, opt)
		if err != nil {
			return result, err
		}

		for _, d := range droplets {
			result[index] = d
			index++
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return result, err
		}

		// set the page we want for the next request
		opt.Page = page + 1

	}

	return result[0:index], nil
}
