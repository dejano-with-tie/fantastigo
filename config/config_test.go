package config

import (
	"os"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	type args struct {
		lc LoadConfig
	}

	tests := []struct {
		name    string
		args    args
		before  func() error
		after   func() error
		want    Config
		wantErr bool
	}{
		{
			name: "Fail_on_missing_env_file",
			args: args{lc: LoadConfig{
				Env:      "missing",
				DirPath:  "fixtures",
				FileName: "test",
			}},
			want:    Config{},
			wantErr: true,
		},
		{
			name: "Config_is_constructed_from_default_env_file",
			args: args{LoadConfig{
				Env:      "",
				DirPath:  "fixtures",
				FileName: "test",
			}},
			want: Config{EnvName: "default", Server: Server{
				Addr:        ":8080",
				Timeout:     3,
				Debug:       false,
				SwaggerUi:   true,
				OpenapiSpec: []string{"spec.yaml"},
			}},
		},
		{
			name: "Config_from_env_specific_file_overrides_default_config",
			args: args{LoadConfig{
				Env:      "dev",
				DirPath:  "fixtures",
				FileName: "test",
			}},
			want: Config{EnvName: "dev", Server: Server{
				Addr:        ":8081",
				Timeout:     3,
				Debug:       false,
				SwaggerUi:   true,
				OpenapiSpec: []string{"spec.yaml"},
			}},
		},
		{
			name: "Env_variable_has_highest_priority_and_overrides_other_configs",
			args: args{LoadConfig{
				Env:      "dev",
				DirPath:  "fixtures",
				FileName: "test",
			}},
			before: func() error {
				eVars := map[string]string{
					"FGO_SERVER_ADDR":        ":8082",
					"FGO_SERVER_SWAGGER__UI": "false",
				}
				return setEnv(eVars)
			},
			after: func() error {
				eVars := []string{
					"FGO_SERVER_ADDR",
					"FGO_SERVER_SWAGGER__UI",
				}
				return unsetEnv(eVars)
			},
			want: Config{EnvName: "dev", Server: Server{
				Addr:        ":8082",
				Timeout:     3,
				Debug:       false,
				SwaggerUi:   false,
				OpenapiSpec: []string{"spec.yaml"},
			}},
		},
		{
			name: "Env_variables_without_expected_prefix_FGO_are_ignored",
			args: args{LoadConfig{
				Env:      "dev",
				DirPath:  "fixtures",
				FileName: "test",
			}},
			before: func() error {
				eVars := map[string]string{
					"GO_SERVER_ADDR":        ":8082",
					"GO_SERVER_SWAGGER__UI": "false",
				}
				return setEnv(eVars)
			},
			after: func() error {
				eVars := []string{
					"GO_SERVER_ADDR",
					"GO_SERVER_SWAGGER__UI",
				}
				return unsetEnv(eVars)
			},
			want: Config{EnvName: "dev", Server: Server{
				Addr:        ":8081",
				Timeout:     3,
				Debug:       false,
				SwaggerUi:   true,
				OpenapiSpec: []string{"spec.yaml"},
			}},
		},
	}

	for _, tt := range tests {

		mustExec(t, tt.before)

		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.lc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})

		mustExec(t, tt.after)
	}
}

func unsetEnv(eVars []string) error {
	for _, k := range eVars {
		if err := os.Unsetenv(k); err != nil {
			return err
		}
	}
	return nil
}

func setEnv(eVars map[string]string) error {
	for k, v := range eVars {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	return nil
}

func mustExec(t *testing.T, fn func() error) {
	if fn == nil {
		return
	}

	err := fn()
	if err != nil {
		t.Errorf("Failed to execute fn: %v", err)
	}
}
