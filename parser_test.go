package halvector

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/koykov/fastconv"
)

type stage struct {
	hal    string
	expect string
	err    error
}

func loadTestDS() (stages []stage, err error) {
	files, err := filepath.Glob("testdata/*.hal.txt")
	if err != nil {
		return nil, err
	}

	fileExists := func(path string) bool {
		_, err := os.Stat(path)
		return err == nil
	}
	fileLoad := func(path string) (contents []byte, err error) {
		if contents, err = os.ReadFile(path); err != nil {
			return
		}
		if len(contents) > 0 && contents[len(contents)-1] == '\n' {
			contents = contents[:len(contents)-1]
		}
		return
	}

	stages = make([]stage, 0, len(files))
	for i := 0; i < len(files); i++ {
		var rawHAL, rawFmt, rawErr []byte
		if rawHAL, err = fileLoad(files[i]); err != nil {
			break
		}
		if ffmt := strings.Replace(files[i], "hal.txt", "fmt.json", 1); fileExists(ffmt) {
			if rawFmt, err = fileLoad(ffmt); err != nil {
				break
			}
		}
		if ferr := strings.Replace(files[i], "hal.txt", "err.txt", 1); fileExists(ferr) {
			if rawErr, err = fileLoad(ferr); err != nil {
				break
			}
		}
		stg := stage{
			hal:    string(rawHAL),
			expect: string(rawFmt),
			err:    errors.New(string(rawErr)),
		}
		stages = append(stages, stg)
	}
	return
}

func TestParser(t *testing.T) {
	stages, err := loadTestDS()
	if err != nil {
		t.Fatal(err)
	}
	for _, stg := range stages {
		t.Run(stg.hal, func(t *testing.T) {
			var buf bytes.Buffer
			vec := Acquire()
			if err := vec.ParseStr(stg.hal); err != nil {
				if stg.err != nil {
					if stg.err.Error() != err.Error() {
						t.Error(err)
					}
					return
				}
				t.Error(err)
				return
			}
			_ = vec.Sort().Beautify(&buf)
			if stg.expect != buf.String() {
				t.Errorf("expect: %s\ngot: %s", stg.expect, buf.String())
			}
		})
	}
}

func BenchmarkParser(b *testing.B) {
	stages, err := loadTestDS()
	if err != nil {
		b.Fatal(err)
	}
	for _, stg := range stages {
		b.Run(stg.hal, func(b *testing.B) {
			b.ReportAllocs()
			var buf bytes.Buffer
			for i := 0; i < b.N; i++ {
				buf.Reset()
				vec := Acquire()
				if err := vec.ParseStr(stg.hal); err != nil {
					b.Error(err)
					return
				}
				_ = vec.Sort().Beautify(&buf)
				exp := buf.Bytes()
				if !bytes.Equal(fastconv.S2B(stg.expect), exp) {
					b.Errorf("expect: %s\ngot: %s", stg.expect, buf.String())
				}
				Release(vec)
			}
		})
	}
}

func BenchmarkLimit(b *testing.B) {
	longHAL := []byte(`pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7,af;q=0.6,sq;q=0.5,de;q=0.4,de-DE;q=0.3,de-AT;q=0.2,de-LI;q=0.1,de-CH;q=0.1,am;q=0.1,ar;q=0.1,an;q=0.1,hy;q=0.1,ast;q=0.1,az;q=0.1,eu;q=0.1,bn;q=0.1,be;q=0.1,nb;q=0.1,my;q=0.1,bs;q=0.1,br;q=0.1,bg;q=0.1,kn;q=0.1,ca;q=0.1,kk;q=0.1,ceb;q=0.1,chr;q=0.1,zh;q=0.1,zh-HK;q=0.1,zh-CN;q=0.1,zh-TW;q=0.1,si;q=0.1,ko;q=0.1,co;q=0.1,hr;q=0.1,ku;q=0.1,ckb;q=0.1,da;q=0.1,sk;q=0.1,sl;q=0.1,es;q=0.1,es-419;q=0.1,es-AR;q=0.1,es-CL;q=0.1,es-CO;q=0.1,es-CR;q=0.1,es-ES;q=0.1,es-US;q=0.1,es-HN;q=0.1,es-MX;q=0.1,es-PE;q=0.1,es-UY;q=0.1,es-VE;q=0.1,eo;q=0.1,et;q=0.1,fo;q=0.1,fil;q=0.1,fi;q=0.1,fr;q=0.1,fr-CA;q=0.1,fr-FR;q=0.1,fr-CH;q=0.1,fy;q=0.1,gd;q=0.1,gl;q=0.1,cy;q=0.1,el;q=0.1,ka;q=0.1,gn;q=0.1,gu;q=0.1,ht;q=0.1,ha;q=0.1,haw;q=0.1,he;q=0.1,hi;q=0.1,hmn;q=0.1,nl;q=0.1,hu;q=0.1,ig;q=0.1,yi;q=0.1,id;q=0.1,en-ZA;q=0.1,en-AU;q=0.1,en-CA;q=0.1,en-IN;q=0.1,en-NZ;q=0.1,en-GB-oxendict;q=0.1,en-GB;q=0.1,ia;q=0.1,ga;q=0.1,yo;q=0.1,is;q=0.1,it;q=0.1,it-IT;q=0.1,it-CH;q=0.1,ja;q=0.1,jv;q=0.1,km;q=0.1,lo;q=0.1,la;q=0.1,lv;q=0.1,ln;q=0.1,lt;q=0.1,lb;q=0.1,mk;q=0.1,ml;q=0.1,ms;q=0.1,mg;q=0.1,mt;q=0.1,mi;q=0.1,mr;q=0.1,mn;q=0.1,ne;q=0.1,ny;q=0.1,no;q=0.1,nn;q=0.1,oc;q=0.1,or;q=0.1,om;q=0.1,pa;q=0.1,ps;q=0.1,fa;q=0.1,pl;q=0.1,pt-PT;q=0.1,qu;q=0.1,rw;q=0.1,ky;q=0.1,rm;q=0.1,ro;q=0.1,mo;q=0.1,ru;q=0.1,sm;q=0.1,sr;q=0.1,sh;q=0.1,sd;q=0.1,so;q=0.1,st;q=0.1,sw;q=0.1,sv;q=0.1,su;q=0.1,tg;q=0.1,th;q=0.1,ta;q=0.1,tt;q=0.1,cs;q=0.1,te;q=0.1,ti;q=0.1,to;q=0.1,tn;q=0.1,tr;q=0.1,tk;q=0.1,tw;q=0.1,uk;q=0.1,ug;q=0.1,wo;q=0.1,ur;q=0.1,uz;q=0.1,wa;q=0.1,vi;q=0.1,xh;q=0.1,sn;q=0.1,zu;q=0.1`)
	fn := func(b *testing.B, limit int) {
		b.ReportAllocs()
		vec := NewVector()
		for i := 0; i < b.N; i++ {
			_ = vec.SetLimit(limit).Parse(longHAL)
			vec.Reset()
		}
	}
	b.Run("no limit", func(b *testing.B) { fn(b, 0) })
	b.Run("limit 50", func(b *testing.B) { fn(b, 50) })
	b.Run("limit 10", func(b *testing.B) { fn(b, 10) })
	b.Run("limit 5", func(b *testing.B) { fn(b, 5) })
	b.Run("limit 3", func(b *testing.B) { fn(b, 3) })
}
