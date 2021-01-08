package proto

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var _proto_micro_mall_users_proto_users_users_swagger_json = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5d\x6d\x8f\xdb\xc6\xf1\x7f\x7f\x9f\x82\xd0\xff\xff\x32\xcd\xa5\x6e\xd1\x17\x7e\x47\x4b\x3c\x5b\xb0\x1e\x0e\x92\xce\xee\xa1\x30\x08\x1e\xb9\xd2\x6d\xcc\x07\x79\x77\x79\xb6\x10\x18\xb0\x8b\xc6\x4d\xd0\xa6\x0e\x10\xc7\x6d\x10\x37\x71\xd0\xc6\x0d\x50\xc0\x0e\x10\xb7\x4e\xce\x75\xfd\x65\x2c\xf9\xfc\x2d\x0a\x92\xa2\x8e\x4b\x2e\x9f\xa9\x3b\xd2\xe0\x02\x09\x7c\x5c\xee\xec\xcc\xec\xfc\x66\x77\x67\x77\xa8\x0f\x36\x38\xae\x81\xaf\x4b\x93\x09\x40\x8d\xb3\x5c\xe3\xcc\xbb\xef\x35\xde\xb1\x9e\x41\x7d\x6c\x34\xce\x72\x56\x3d\xc7\x35\x08\x24\x2a\xb0\xea\xa7\xc8\x20\xc6\xa6\x06\x65\x64\x88\x9a\xa4\xaa\xa2\x89\x01\xc2\xa2\xf3\xd8\xfe\xb7\xf3\xff\x77\xed\x27\x36\x29\x8e\x6b\x1c\x00\x84\xa1\xa1\x5b\x04\x96\xff\xe4\x74\x83\x70\x18\x90\xc6\x06\xc7\xdd\xb4\x3b\x94\x0d\x1d\x9b\x1a\xc0\x8d\xb3\xdc\x6f\x9c\x56\xd2\x74\xaa\x42\x59\x22\xd0\xd0\x37\xdf\xc7\x86\x6e\xbd\x7b\xc5\x7e\x77\x8a\x0c\xc5\x94\x13\xbe\x2b\x91\x7d\x7c\x2c\xc9\xe6\xc1\xcf\x6d\x0e\x37\x25\x59\x36\x4c\x9d\x6c\xca\xfb\x12\x9a\x80\xd5\x0b\x56\x0b\x93\x78\xfe\xb4\xf4\x63\x6a\x9a\x84\x66\x16\xfb\xaf\xef\x7d\xb7\xf8\xe8\xd9\xd1\xd3\x47\x8b\x8f\x9e\xcd\x3f\xfc\x70\x7e\xeb\x3f\x4b\x19\xed\x17\x8d\x29\x40\x36\x0f\x6d\xc5\x7a\x79\xc7\xd2\xc4\x10\xa0\x03\x28\x03\xd1\xfa\x83\x77\xfa\x6c\x3a\x5d\x7a\x1a\x22\x80\xa7\x86\x8e\x01\xa6\x3a\xe6\xb8\xc6\x99\xf7\xde\xf3\x3d\xe2\xb8\x86\x02\xb0\x8c\xe0\x94\x2c\x55\xca\x73\xd8\x94\x65\x80\xf1\xd8\x54\x39\x97\xd2\xbb\x1e\xf2\x8e\x10\xf2\x3e\xd0\xa4\x00\x31\x8e\x6b\xfc\x3f\x02\x63\x8b\xce\xff\x6d\x2a\x60\x0c\x75\x68\xd1\x5d\x0e\x63\x80\xe9\xc1\x92\x7c\x83\x22\x72\xd3\xf3\xd7\x4d\x6f\xbf\x0d\x05\x8c\x25\x53\x25\xf1\x32\xe8\x9c\xa9\x83\x1b\x53\x20\x13\xa0\x70\x00\x21\x03\x15\x27\x0a\x32\x75\x02\x35\x20\x58\x54\x23\x18\xdf\x60\x88\xd0\x98\x4a\x48\xd2\x00\x01\xe8\xd8\xd8\x9c\xe2\x93\x47\x97\x34\x1b\x20\x7b\x86\x32\xf3\xf3\x0b\xf5\xb0\x1a\x04\xae\x99\x10\x01\xcb\x5a\x08\x32\xc1\x9a\x86\xec\x9a\x09\x30\x49\x22\xf8\x15\x8f\xe0\x44\x9a\xf8\x45\xa6\x0c\xfa\x98\xde\x95\x0d\x2f\x9d\xa5\xf2\x8e\x81\x06\x6e\x40\x4c\x36\xa7\xfb\x86\x4e\xa1\x6c\x02\xc2\x51\xb6\xf8\xdb\xad\xc5\xd7\xdf\x3a\x58\x5b\xfc\xe5\xc9\xfc\xd3\x47\x8b\x1f\xbe\x9b\xdf\xf9\xe3\xcf\x16\x1f\xff\x61\xf1\xe0\x70\x7e\xf7\x59\x52\xd4\x35\xf7\x81\x7c\xd5\x7a\x72\x6e\xb6\x6d\x73\x50\x01\xd0\xf9\x79\xae\x31\xe7\x96\x10\xcc\xd9\xc6\x8e\x66\xa2\x6c\x28\x80\x8d\xbd\x6b\x26\x40\x51\xe0\x1b\x4b\x2a\xf6\xa3\x8f\xcc\xa6\x36\x75\x4c\x10\xd4\x27\x8d\x30\x75\x87\xf0\x34\xf5\x59\xdb\xfa\x98\x59\x2f\x7a\xa1\x02\x74\x02\xc9\x2c\x0b\x74\x8f\x0e\xff\xf9\xea\xf9\x8b\xd4\x50\x6d\xbb\x5d\x56\x09\xab\x2e\xd3\x35\x58\xdd\x52\x83\x35\x23\x33\xd7\x95\xb7\xc1\x6f\x78\xf7\x0f\x31\x3e\xe3\xe8\x4f\xcf\xe6\x77\xef\x3b\x3e\xe3\xd5\xcb\x6f\x16\xb7\x9f\x24\xf5\x19\xe7\x01\xb1\xc1\x67\x75\x56\x01\x6f\xe1\x61\xb7\xf6\x13\x6e\x09\x81\x81\x09\xd7\x07\x03\x5f\xed\xd8\x40\x9a\x64\xa9\xb6\x01\x75\xf2\xab\x5f\x9e\x30\x4a\xdc\x9d\xa8\x08\x95\x2c\x80\x59\xb6\x6e\x2b\x29\x31\xc3\xb3\xda\x95\x1c\x38\x2b\x9e\x6b\xf4\xb8\x25\x1c\x3d\xa2\x0a\x31\x29\x1e\x42\x12\x42\x52\x70\x87\x4b\x80\xe6\xb7\x19\x2e\x16\x78\xd1\xd0\xf3\x8d\xa1\x1d\x22\x52\x55\x20\x5b\xca\xdc\x5a\x35\xd3\x4c\x95\xc0\x93\x46\xac\x6c\xad\xfa\x32\x6d\x66\xef\xfd\xfb\xe8\xe9\x3f\xe6\x0f\xbe\x9f\xff\xf5\x56\xea\x75\x71\x0b\xa8\xf0\x00\xa0\x59\x55\x66\x3b\x26\xe3\x35\x72\xdd\x52\xee\x79\x2f\xc1\x42\x55\x59\x0e\xab\x08\x15\x5c\x16\x3f\x03\x75\x02\x26\x00\xc5\x38\x9a\x5f\x9c\xa9\x88\xa3\x71\x55\x9c\x65\x61\xe0\xf8\x9a\x4c\xeb\xe9\xaa\x79\x1a\x06\xdb\xb5\x9f\x71\x4b\xe5\xfd\x8c\x35\xc4\xa2\xc7\xd9\x14\xcf\x36\xdb\x6b\x44\xf9\x8c\xc2\xc0\xef\xb6\x8e\x3c\x75\x7a\xf5\xf2\xf1\xe2\xde\x4f\xd9\x61\xdd\x35\x14\x38\x9e\x55\x11\xd9\x6c\xce\x6b\x70\xbb\xa5\xac\xa7\x50\x61\xe3\x76\xba\x47\x51\xf6\x9c\x3a\x86\x7a\xe2\x8d\xf6\xe2\xeb\x6f\x17\x1f\xff\x37\x4b\x64\x6a\x0b\xea\x4a\x95\x42\x53\x5e\x7e\x6b\x78\xb9\xa5\xde\x5d\x9f\xc2\xa2\x17\xea\x07\x90\x00\xe7\x80\x20\x29\x4e\x1f\xfe\xb4\xf8\xe4\xf1\x9b\xdf\xde\x3a\x7a\xf2\xec\xf5\xc3\xdb\xce\x32\xd8\xf9\xf3\xd5\xe1\x61\x86\x68\xf2\xb9\x59\xdb\xe6\xa2\x49\x9f\x52\x94\x16\xbe\x4b\xd6\xbd\x6c\xd7\x28\x76\x4b\x08\x8a\xbd\x76\xf6\x76\x1c\xb8\xa4\xbb\x65\x11\x3c\x76\x49\x7f\xbf\x82\x82\x4c\x65\x6e\x58\x04\xb9\xae\xc1\xe2\x96\xfa\xd8\xf6\xc4\xa1\xab\x1a\x13\x9b\xeb\x64\x37\x10\xbf\x78\x3e\x7f\xf1\x79\x52\x7c\x76\x2c\xd2\xd6\x93\x2a\xc0\x72\xc5\x6c\x8d\x46\xb7\x94\x75\x7f\xe7\x19\xaa\xd3\xdd\xd2\x69\x00\xc9\xfb\x92\x4e\xf0\xa6\x84\xb1\x21\x43\x89\x80\x4d\xbc\x6f\x4c\xc5\x3d\x13\x43\x1d\x60\x9c\x14\x57\xf3\xcf\xef\xd8\x77\x7a\x7f\x38\xba\x7d\x6f\x7e\xf8\xe7\x37\x9f\x45\x2d\x1b\xbb\x6e\xaf\xab\x08\x8b\xfb\x80\x77\xb9\x18\xee\x1b\xd3\x2a\x80\x8e\xcd\x79\x8d\x40\xb7\x94\x15\x81\x61\xe3\x56\x18\x1c\xfd\x36\x9e\x05\x92\xa6\x02\x49\x62\x00\x3e\xfe\x66\xf1\xf0\x47\x07\x86\xef\x1c\xbd\xbc\x37\xff\xf2\x2b\xe7\xd1\xeb\xe7\x9f\x2d\xbe\x7a\x90\x09\x8d\x5d\x89\x00\x04\x25\x95\xb7\x19\xa9\x12\x1a\x29\xce\x6b\x34\xba\xa5\xf4\x68\xf4\x8d\x5b\x99\xd0\xa8\x2d\x59\x4b\xb7\x3f\x74\x00\x79\xf4\xaf\xdf\x2d\xee\x7f\x91\x06\x83\xe7\x01\x09\x28\xa5\x0a\x08\x64\xf1\x5d\xe3\xcf\x2d\x21\xf8\x73\x4d\x6b\x2d\xa7\x73\x6b\xbc\xb4\x97\x00\x55\xc7\x07\x74\x06\x8e\x5b\x3d\x2e\xee\x7e\xfa\xea\xf0\xef\xe9\xb1\x52\x49\xa0\xd4\x28\xa9\xfe\x2c\x55\x86\x09\x6a\x2a\x61\x7c\xdd\x40\xca\x26\x02\x18\x24\x5e\x2d\xbe\xf9\xfd\x27\xaf\x5f\x3c\x76\x82\x21\xf3\x27\x77\x5e\x3f\xbc\x9d\x34\x24\xb2\xbd\xec\x6f\x60\x77\x57\x01\xa0\x51\x0c\xd7\x20\x73\x4b\x59\x41\xe6\x1b\xae\xd3\x0d\x8f\x20\x30\x81\x98\xd8\xb9\xdd\x1f\x24\x9a\xc8\x9c\x44\x4b\xe7\x50\x20\x29\xa2\x06\x6e\x27\x15\x00\x93\xcb\x6b\x8d\x23\xb7\x94\x15\x47\xc7\x23\x75\xba\x10\xc2\x44\x22\x20\xfb\x95\xef\xd4\xd7\xbc\x87\x56\x7f\x55\x80\x12\xcd\x71\x0d\x28\xb7\xd4\x97\x46\x4e\xf6\xd2\x88\x03\xd0\x54\x87\x68\x8b\x2f\x9f\x2e\xee\x7f\xef\x00\xd4\x59\x43\x2e\x6e\x25\x5e\x40\xee\x4c\x15\x89\x00\xeb\x91\x7d\x0a\x52\x19\xbc\xb2\xf8\xae\x51\xeb\x96\xb2\x4e\x83\xec\x51\x3b\xd9\x29\x71\xf5\x91\x1e\x0f\x77\x2b\x31\x9c\x6f\x02\xed\x99\x63\x5e\xa7\xf2\x14\x5c\x3f\x63\xec\xbd\x0f\xe4\x63\xaf\x67\xbd\x3e\x05\x88\x40\x1f\x3e\xec\xf7\x45\x13\xa9\x7e\xd4\x84\x1d\xcb\x7b\x47\xf5\x40\x52\x4d\x10\xd3\x90\xb2\xe3\x63\x27\xb7\x37\x23\x1e\xc1\x6f\x32\x5d\x0d\x65\x70\x39\x44\x04\x3e\x02\x89\xe5\xf3\xdd\x85\xa3\xda\xb1\x2e\xaf\x87\x5f\x5d\xf7\x52\xd5\x00\xc6\xd2\x24\x4e\x6f\xcc\xa6\x0a\x20\x12\x54\x03\x1e\x2e\x7c\x46\x0a\x99\x8f\x42\x2c\xdf\x6b\x52\x6c\x48\x33\xc7\xc9\x06\xcc\x32\x49\x74\xe4\xb0\x12\x18\x2b\x9f\x35\x34\x80\x6e\x6a\x14\x28\x1a\xdb\xc2\x60\xd8\xef\xf1\x1d\xaf\x4f\x6f\xf6\xbb\xdb\x7c\x6f\xd7\xfb\x68\xb8\x3b\x1c\x09\x5d\x97\xbd\x15\xce\x3c\x2e\xf2\x98\x52\x90\xcb\xd0\x8f\xbc\xe4\x30\x2f\xd9\xd0\x34\x43\xf7\x8f\x49\xd4\x0a\xce\x6e\x10\x9c\x03\xbc\x23\x0d\xb1\x68\x7f\xc1\x27\x6c\xa8\xf7\x0c\x43\x05\x92\x1e\x07\xa1\x04\x19\x81\xa5\x90\x3c\x01\xf3\x81\x4f\x7d\x54\x85\x71\x7a\xc6\x2f\x3b\xd7\x06\xd4\xb3\x63\xb8\x49\x63\x75\x67\xd8\x8a\x04\xaa\xf5\x3a\x8b\x05\x8a\xeb\x5c\x0a\x0b\xfa\xef\xc8\xad\x2f\xb1\x2f\x1e\xb3\xbd\x36\x9e\x24\xf6\xd8\xe1\xea\x15\x34\x09\xaa\xf9\x44\xd2\x09\xd0\x43\xdd\x42\x0a\x56\x98\xa9\x07\xa5\xb0\x4e\xaf\x1b\xd4\xc7\x86\xb3\x7d\x5b\xd3\x94\xb7\xfa\x78\x9b\xa5\x86\xae\x04\xf5\xd4\x13\x5f\xe4\xc9\x65\x19\xf5\x99\x82\x64\x40\x30\x3b\xaf\x26\xa1\x56\xd8\x9f\x8e\x28\xa3\x46\x4e\xc6\xc2\x5c\x5d\x50\x2a\x4c\x6e\x64\xa1\x89\x06\x65\xd4\x68\x0a\x92\xab\x74\xad\x14\x7a\x28\xef\x42\x26\x46\x0f\x85\x1b\x15\x95\x59\x9a\xd1\xaa\x58\x17\xf2\xab\xa0\xcc\x82\x8d\xaa\x72\xc6\x54\x94\xfc\x6d\xdc\x72\xd6\x64\xd9\x17\x7e\x2d\x61\x8b\xdf\xe9\x8c\xc4\xd1\xee\xb6\x20\x6e\xf1\x9d\xa1\xe0\x5d\x07\x52\xb5\xa3\xc1\x8e\x10\xb9\x2a\x64\xd0\x0a\xf2\x6c\x47\x65\xce\xcd\xdc\x53\xbf\x3c\xe3\x65\x87\x30\xc5\xab\x74\x7e\x6a\x9c\x82\xed\xfe\xb7\xaf\x2b\x17\xad\x66\xcc\x11\x03\xbe\xa5\x5e\x1c\x45\x67\x69\xc8\x24\xe5\xcf\x2f\x8a\x23\xd5\x35\xf6\xa0\x0a\x9c\xbc\x20\x36\xc1\xeb\x01\x59\x33\xac\x20\x29\x1d\x64\x32\x9b\x6e\xff\x5c\xbb\x23\x88\xdb\x17\xfa\x3d\xca\x60\x84\x2e\xdf\xee\x44\x1a\x09\xd5\x32\x84\xb5\x91\x71\x15\xe8\xb9\x62\x63\x3e\x02\xb9\x14\x95\x23\x32\x72\xb9\xe5\x55\xce\x25\x61\xd0\xde\xda\x15\x9b\xfd\x16\xa5\xb3\x51\xff\xa2\xd0\x8b\x8e\x8b\x5c\x6e\x85\xa9\xca\x9b\x53\x90\x1b\x4a\x84\x96\x34\xce\x5a\x8f\xb5\x93\xd4\x56\xe3\x88\x79\xfc\x02\x93\xe4\x01\x40\x70\x3c\xf3\xa7\xba\x26\x22\x7d\xc9\x6e\x1a\xbe\x4d\x64\x9a\x4c\xac\xf0\x76\xa3\xc4\x96\x44\xe5\xea\x94\x6e\xa2\x5a\x86\x69\xc4\x42\xc1\xe3\x51\x7b\x0e\x81\x0b\xf7\xa3\x11\x86\x94\x41\xd2\xe8\xdc\x82\x1c\x72\x7b\xef\x89\x66\x8a\xd6\xfb\x8e\x24\xbd\x2a\xb0\xb3\x7e\xd6\x40\x17\x62\x71\x95\x5a\x54\x44\x10\x34\x26\xe1\xa6\x14\x30\x4a\xc0\x3d\xf3\xa2\x7b\x49\x4d\xc3\xce\x43\x11\xcd\x35\x50\xb6\x54\x5e\x4c\x14\x2c\x26\xf3\xa3\x14\x66\xe1\x45\x1b\x89\x80\x43\xf6\x23\xa9\xb4\xfa\xf4\x34\xd5\xc1\x0d\x22\xae\x6f\xa4\x21\x16\x8d\x03\x10\x7a\x7a\x97\xc9\x03\x50\xc1\xac\x1c\x23\xbc\x0e\x79\xd7\x89\x47\xf7\x56\xa2\x28\x29\x4a\xa6\xf3\xd0\x7d\x20\xa9\x64\x5f\x94\x25\xa4\x88\x7a\x68\x6c\x23\x8a\x02\xe3\xbb\xf6\x5c\x31\x26\xbc\x26\x64\x10\xe9\x46\x1e\x79\x65\x04\x24\x02\x44\x02\xb5\x4c\xe7\xbd\xa6\x7d\x0b\x21\x5d\xf3\x14\x00\x28\x60\xfe\x28\x38\xa0\xec\x91\x7d\x75\x17\x28\xed\x7e\xa2\xef\x36\xf4\xed\x29\x52\x29\xa6\xa4\xfe\xbf\x70\xf7\x10\xa1\x14\xcf\xe2\x37\x97\x1a\x3c\x5f\xb0\xc8\x00\x01\xe6\x9a\x3d\x8b\xf1\x47\x7e\x1a\x2d\x87\x80\x85\xdb\x69\xbe\xa8\x1f\x3b\x20\x4c\x39\x95\x93\xb0\x9d\xc8\xef\x07\x96\x02\x55\xe1\xdc\xaf\x06\xa6\x0b\x88\x94\xcf\x34\xd6\xb2\x26\xb2\x36\x5c\xd3\x2c\x58\x32\xa6\xe2\x54\x95\x88\xd5\x49\xc6\xe6\x0a\xb0\x6f\xac\xe5\x87\x23\x6d\xfc\xd9\x2e\x1b\x0c\x04\x7e\x44\x85\xc0\x76\xb6\x5b\xbe\x27\x2d\xa1\x23\xd0\x4f\xf8\x9d\x56\x7b\x14\x7d\x2b\xc1\xa1\x1b\xe4\x99\x99\x62\x52\xb2\xf5\x63\x31\x91\x5d\x76\xee\x53\xc9\x61\xeb\xcf\x5c\xc8\x33\x32\x18\x20\x71\x79\x17\x35\x35\x4c\x30\xb8\xb1\x8e\x1d\xda\x49\x4e\xa3\x71\xc7\x18\x89\xd4\xb0\x9e\x05\x39\x54\x72\xed\x3f\xec\xe3\x73\x64\xff\x7e\xe6\x1a\xa2\x12\x3a\x91\x64\x92\x79\x7b\x15\x71\x29\x34\xbb\xbe\xa6\xc1\x13\xb2\x68\x86\x92\x00\xac\x4c\x0e\x81\xda\xe0\xe2\xe0\x2d\xf9\x84\x69\x69\x56\xc3\x94\x4a\xa0\xfb\xca\xb0\x75\x62\x7d\x70\x33\xfb\xe0\x90\xb0\xf0\x78\xfc\x54\x3a\xdc\x69\x36\x85\xe1\x90\x3a\x82\x1b\x0c\xfa\x03\xfa\x32\x9f\x30\x10\x7b\xfd\x91\x28\xfc\xba\x3d\x1c\x05\x6a\x02\x4f\xbb\xc2\xa0\x79\x81\xef\x8d\xd8\x6d\x56\xb5\x81\x1a\xbe\xd9\xec\xef\x44\x55\x84\xf3\xb0\x7d\xb9\x65\xd7\x76\xf9\x51\xf3\x42\xa0\xb6\xd3\x3f\xdf\xee\xd9\xf5\x7c\xa7\xd3\xbf\x1c\xa8\x6f\x09\x9d\xf6\x25\x61\xb0\x2b\xb6\x7b\x5b\xfd\x90\x1e\xe8\x77\x98\x9c\x8c\x06\x7c\x6f\xc8\x37\x47\xed\x7e\x4f\xdc\xe2\xdb\x1d\xa1\xc5\x12\xa2\xd3\x6f\x5e\x64\x3d\x6f\xf7\x2e\xf1\x9d\x76\x2b\xd0\x71\xf3\x02\x3f\x38\x2f\x88\x83\x9d\x5e\x58\x15\x63\x0c\xbd\xd5\xa3\x01\xdf\x12\xc4\x5e\x5f\x14\xba\xdb\xa3\xdd\xc8\x05\x90\x4b\x29\x68\x64\x51\x59\x11\x25\x5b\x08\x31\xa7\x9f\x98\xbd\x93\x27\xb1\x28\x09\xe4\x22\x33\x7b\x4a\xe1\x1b\x23\x78\x0f\xfb\xbd\xda\x7c\xc3\x58\xd0\x4d\xba\x02\xbe\x41\x42\x4d\xac\xcb\x9f\xfa\x4a\xb9\x41\xf7\xa6\x34\x84\x4c\xf8\xe9\x4f\xbb\x57\x17\xac\xd9\x9c\x6a\x56\x8f\x99\x36\x67\x26\x11\x09\x92\x14\x90\x71\x65\x64\x4c\x45\x8d\xde\xee\xc6\x49\x42\x6f\x93\x33\x1a\x5d\xe5\xd0\x52\xc6\x43\x03\xe6\x4f\xd9\x45\x8f\x7c\xb4\xac\x54\x10\x29\xcf\x12\xa7\x78\x59\x57\xbf\x48\x92\x60\x09\xcf\x24\xa0\xd9\x71\x4d\x31\xf3\x86\x48\x42\x20\x80\x92\xe4\xb9\x5b\x40\x11\xb3\x12\x50\xa5\x3d\x10\xba\x13\xcb\xec\x5b\x63\x9d\x27\xc4\x22\x3b\xfd\x34\x0a\x63\xf4\x95\xc2\xa4\x66\xf7\x96\x60\xcb\x1b\x6b\xad\x03\x09\xee\x3b\x6f\x65\x20\xa1\xf8\xf3\xee\x65\x80\x22\xeb\xf0\x94\x32\x06\x51\xaa\x53\x50\x2a\xd9\xa7\x64\xce\x26\xab\xa7\x38\x35\x50\xaf\xc5\x5c\x2c\xdb\xa5\x3f\x21\x9d\x6c\xac\x57\x43\x46\x20\x51\xed\x57\x97\xbf\xb2\xf0\xe3\xf3\xa3\x47\xb7\x97\x3f\x24\xc4\x36\x08\xcf\xce\x2f\xd7\x52\x34\x5d\x72\x9e\x47\xe6\x04\xc6\x9d\xe9\x54\x6e\xc3\xfa\xef\xe6\xc6\xff\x02\x00\x00\xff\xff\xc9\x32\x65\x99\x71\x88\x00\x00")

func proto_micro_mall_users_proto_users_users_swagger_json() ([]byte, error) {
	return bindata_read(
		_proto_micro_mall_users_proto_users_users_swagger_json,
		"proto/micro_mall_users_proto/users/users.swagger.json",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"proto/micro_mall_users_proto/users/users.swagger.json": proto_micro_mall_users_proto_users_users_swagger_json,
}
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"proto": &_bintree_t{nil, map[string]*_bintree_t{
		"micro_mall_users_proto": &_bintree_t{nil, map[string]*_bintree_t{
			"users": &_bintree_t{nil, map[string]*_bintree_t{
				"users.swagger.json": &_bintree_t{proto_micro_mall_users_proto_users_users_swagger_json, map[string]*_bintree_t{
				}},
			}},
		}},
	}},
}}
