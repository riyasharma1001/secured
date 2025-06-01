package security

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

type Obfuscator struct {
	minifier *minify.M
}

func init() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())
}

func NewObfuscator() *Obfuscator {
	m := minify.New()
	m.AddFunc("text/javascript", js.Minify)
	return &Obfuscator{minifier: m}
}

func (o *Obfuscator) Minify(code []byte) []byte {
	minified, err := o.minifier.Bytes("text/javascript", code)
	if err != nil {
		return code
	}
	return minified
}

func (o *Obfuscator) Obfuscate(code []byte) []byte {
	jsCode := string(code)

	// 1. Variable name mangling
	jsCode = o.mangleVariables(jsCode)

	// 2. Add dead code
	jsCode = o.insertDeadCode(jsCode)

	// 3. Add code transformations
	jsCode = o.transformCode(jsCode)

	// 4. Add timing checks
	jsCode = o.addTimingChecks(jsCode)

	// 5. Add string encryption
	jsCode = o.encryptStrings(jsCode)

	// 6. Add control flow flattening
	jsCode = o.flattenControlFlow(jsCode)

	// 7. Add self-defending code
	jsCode = o.addSelfDefense(jsCode)

	return []byte(jsCode)
}

func (o *Obfuscator) mangleVariables(code string) string {
	replacements := map[string]string{
		"function": fmt.Sprintf("_0x%x", rand.Int31()),
		"return":   fmt.Sprintf("_0x%x", rand.Int31()),
		"var":      fmt.Sprintf("_0x%x", rand.Int31()),
		"let":      fmt.Sprintf("_0x%x", rand.Int31()),
		"const":    fmt.Sprintf("_0x%x", rand.Int31()),
	}

	for old, new := range replacements {
		code = strings.ReplaceAll(code, old, new)
	}
	return code
}

func (o *Obfuscator) insertDeadCode(code string) string {
	deadCode := []string{
		fmt.Sprintf(`if(false){console.log("%x")}`, rand.Int31()),
		fmt.Sprintf(`while(false){console.log("%x")}`, rand.Int31()),
		fmt.Sprintf(`try{throw "%x"}catch(e){}`, rand.Int31()),
	}

	for _, dc := range deadCode {
		pos := rand.Intn(len(code))
		code = code[:pos] + dc + code[pos:]
	}
	return code
}

func (o *Obfuscator) transformCode(code string) string {
	// Convert to array-based access
	transformed := fmt.Sprintf(`
        (function(){
            const _0x%x = [%s];
            return function(){
                return _0x%x.join('');
            }
        })()()
    `, rand.Int31(), o.splitToArray(code), rand.Int31())
	return transformed
}

func (o *Obfuscator) splitToArray(code string) string {
	var parts []string
	size := 3
	for i := 0; i < len(code); i += size {
		end := i + size
		if end > len(code) {
			end = len(code)
		}
		part := code[i:end]
		parts = append(parts, fmt.Sprintf("'%s'", part))
	}
	return strings.Join(parts, ",")
}

func (o *Obfuscator) addTimingChecks(code string) string {
	wrapper := `
        (function(){
            const start = Date.now();
            const result = (function(){
                %s
            })();
            if(Date.now() - start < 100){
                return result;
            }
            return null;
        })()
    `
	return fmt.Sprintf(wrapper, code)
}

func (o *Obfuscator) encryptStrings(code string) string {
	key := make([]byte, 16)
	// Use crand for cryptographic operations
	crand.Read(key)

	wrapper := `
        (function(){
            const d = function(s){
                return atob(s).split('').map(c => 
                    String.fromCharCode(c.charCodeAt(0) ^ %d)
                ).join('');
            };
            return (function(){
                return d("%s");
            })();
        })()
    `
	encoded := base64.StdEncoding.EncodeToString([]byte(code))
	return fmt.Sprintf(wrapper, key[0], encoded)
}

func (o *Obfuscator) flattenControlFlow(code string) string {
	wrapper := `
        (function(){
            const states = {%s};
            let state = 0;
            while(state !== -1){
                state = states[state]();
            }
        })()
    `

	// Split code into states
	states := o.generateStates(code)
	return fmt.Sprintf(wrapper, states)
}

func (o *Obfuscator) generateStates(code string) string {
	// Simple state machine implementation
	return fmt.Sprintf(`
        0:function(){%s;return -1;}
    `, code)
}

func (o *Obfuscator) addSelfDefense(code string) string {
	wrapper := `
        (function(){
            const _0x%x = function(){
                if(new Error().stack.includes("chrome-extension://"))
                    throw new Error();
                if(window.outerHeight - window.innerHeight > 100)
                    throw new Error();
                if(window.Firebug || window.console.profiles)
                    throw new Error();
                return true;
            };
            if(_0x%x()){
                return (function(){
                    %s
                })();
            }
        })()
    `
	checkName := rand.Int31()
	return fmt.Sprintf(wrapper, checkName, checkName, code)
}
