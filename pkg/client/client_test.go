package client

import (
	"context"
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCreateSecret(t *testing.T) {
	testcases := map[string]struct {
		namespace string
		name      string
		data      map[string][]byte
	}{
		"success": {
			namespace: "test",
			name:      "example",
			data: map[string][]byte{
				"foo": []byte("bar"),
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			secret := &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: tc.name,
				},
				Data: tc.data,
			}

			clientset := fake.NewSimpleClientset()
			client := &clientImpl{
				clientset: clientset,
			}

			ctx := context.Background()

			s, err := client.CreateSecret(ctx, tc.namespace, secret)
			if err != nil {
				t.Errorf("want no error, got %q", err)
			}

			if s.Name != secret.Name {
				t.Errorf("secret name want %q, got %q", s.Name, secret.Name)
			}

			if _, err := clientset.CoreV1().Secrets(tc.namespace).Get(ctx, tc.name, metav1.GetOptions{}); err != nil {
				t.Errorf("secret %s not found, error: %q", tc.name, err)
			}
		})
	}
}

func TestGetSecret(t *testing.T) {
	testcases := map[string]struct {
		namespace string
		name      string
		data      map[string][]byte
	}{
		"success": {
			namespace: "test",
			name:      "example",
			data: map[string][]byte{
				"foo": []byte("bar"),
			},
		},
		// TODO: error case here
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			secret := &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      tc.name,
					Namespace: tc.namespace,
				},
				Data: tc.data,
			}

			clientset := fake.NewSimpleClientset(secret)
			client := &clientImpl{
				clientset: clientset,
			}

			s, err := client.GetSecret(context.Background(), tc.namespace, tc.name)
			if err != nil {
				t.Errorf("want no error, got %q", err)
			}

			if s.Name != secret.Name {
				t.Errorf("secret name want %q, got %q", s.Name, secret.Name)
			}
		})
	}
}

func TestListSecrets(t *testing.T) {
	testcases := map[string]struct {
		namespace string
		secrets   []runtime.Object
	}{
		"success": {
			namespace: "test",
			secrets: []runtime.Object{
				&v1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example1",
						Namespace: "test",
					},
					Data: map[string][]byte{
						"foo": []byte("bar"),
					},
				},
				&v1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "example2",
						Namespace: "test",
					},
					Data: map[string][]byte{
						"baz": []byte("qux"),
					},
				},
			},
		},
		// TODO: error case here
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset(tc.secrets...)
			client := &clientImpl{
				clientset: clientset,
			}

			ss, err := client.ListSecrets(context.Background(), tc.namespace)
			if err != nil {
				t.Errorf("want no error, got %q", err)
			}

			if got, want := len(ss.Items), 2; got != want {
				t.Errorf("want %d items, got %d", want, got)
			}
		})
	}
}

func TestUpdateSecret(t *testing.T) {
	testcases := map[string]struct {
		namespace string
		oldSecret *v1.Secret
		newSecret *v1.Secret
	}{
		"success": {
			namespace: "test",
			newSecret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example",
					Namespace: "test",
				},
				Data: map[string][]byte{
					"foo": []byte("bar"),
				},
			},
			oldSecret: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example",
					Namespace: "test",
				},
				Data: map[string][]byte{
					"foo": []byte("qux"),
				},
			},
		},
		// TODO: error case here
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset(tc.oldSecret)
			client := &clientImpl{
				clientset: clientset,
			}

			ctx := context.Background()

			_, err := client.UpdateSecret(ctx, tc.namespace, tc.newSecret)
			if err != nil {
				t.Errorf("want no error, got %q", err)
			}

			s, err := clientset.CoreV1().Secrets(tc.namespace).Get(ctx, tc.oldSecret.Name, metav1.GetOptions{})
			if err != nil {
				t.Errorf("secret %q not found, error: %q", tc.oldSecret.Name, err)
			}

			if !reflect.DeepEqual(s.Data, tc.newSecret.Data) {
				t.Errorf("secret data want %#v, got %#v", tc.newSecret.Data, s.Data)
			}
		})
	}
}
