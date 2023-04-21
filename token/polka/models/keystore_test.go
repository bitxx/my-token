package models

import (
	"mytoken/core/utils/hexutil"
	"testing"
)

func Test_keystore_Sign(t *testing.T) {
	type fields struct {
		Encoded  string
		Encoding *encoding
		Address  string
	}
	type args struct {
		msg      []byte
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "keystore1",
			fields: fields{
				Encoded: "jC9MOH7OPYbHdJtiOWFW0lpMUCFO4nASKjzqHvXpEiYAgAAAAQAAAAgAAACm2Dm/CZ98R1uy34lMj7tr9+i3ERCFoeCSdNwOScsyDkvLwhVGv6qxOzmdiR7vzgRgEizMQbq17k0C1Tk59WyDnf9OfaGQTenQQpnFPiXxcmDa6TXQvF7Eq8VYw009ANLmDTIQ125JdQX6edYY85ZFpLiOltXiad44mhS1mC8OSCcOHsViVrk3Lk0eMsClYS1SUzv3QDCoHChFu6Za",
				Encoding: &encoding{
					Content: []string{"pkcs8", "sr25519"},
					Type:    []string{"scrypt", "xsalsa20-poly1305"},
					Version: "3",
				},
				Address: "5Gc8bR5p9JeCY3dpCvdonRWn79UxhKycDb8aC7xfqQPqWhr8",
			},
			args: args{
				msg:      []byte("123tfyyyufuuyyyyyyyyyyyyyyyyygyghcfgchgfdfsersssss65e766666f66k7fffff7fk6fuf56d65s4d5swaaa33aaaa2a3us5sd"),
				password: "111",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Keystore{
				Encoded:  tt.fields.Encoded,
				Encoding: tt.fields.Encoding,
				Address:  tt.fields.Address,
			}
			got, err := k.Sign(tt.args.msg, tt.args.password)
			if err != nil {
				t.Error(err)
			}
			gotHex := hexutil.HexEncodeToString(got)
			t.Log(gotHex)
			//if (err != nil) != tt.wantErr {
			//	t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
			//	return
			//}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("Sign() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
