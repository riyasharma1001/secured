package security

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	mrand "math/rand" // Alias math/rand to avoid conflicts
	"time"
)

type Pipeline struct {
	encryptor   *Encryptor
	obfuscator  *Obfuscator
	antiDebug   *AntiDebug
	wasmEnabled bool
}

type Config struct {
	EncryptionKey   string
	EnableAntiDebug bool
	EnableWASM      bool
}

func NewPipeline(config *Config) *Pipeline {
	return &Pipeline{
		encryptor:   NewEncryptor([]byte(config.EncryptionKey)),
		obfuscator:  NewObfuscator(),
		antiDebug:   NewAntiDebug(config.EnableAntiDebug),
		wasmEnabled: config.EnableWASM, // Set the WASM flag
	}
}

func (p *Pipeline) Process(code []byte) ([]byte, error) {
	// 1. Initial obfuscation
	obfuscated := p.obfuscator.Obfuscate(code)

	// 2. Add polymorphic code
	polymorphic := p.addPolymorphicLayer(obfuscated)

	// 3. Add control flow flattening
	flattened := p.flattenControlFlow(polymorphic)

	// 4. String encryption
	stringEncrypted := p.encryptStrings(flattened)

	// 5. Add anti-analysis
	protected := p.addAntiAnalysis(stringEncrypted)

	// 6. Apply WASM transformation if enabled
	var processed []byte
	if p.wasmEnabled {
		processed = p.wrapWithWASM(protected)
	} else {
		processed = protected
	}

	// 7. Final encryption
	encrypted, err := p.encryptor.Encrypt(processed)
	if err != nil {
		return nil, err
	}

	// 8. Add advanced loader
	return p.addAdvancedLoader(encrypted), nil
}

// New method for adding anti-analysis
func (p *Pipeline) addAntiAnalysis(code []byte) []byte {
	template := `
    (function(){
        // Performance checks
        const perfStart = performance.now();
        
        // VM Detection
        const checkVM = () => {
            try {
                document.createEvent("TouchEvent");
                return false;
            } catch(e) {
                return true;
            }
        };

        // Browser fingerprinting
        const getBrowserFingerprint = () => {
            return navigator.userAgent + screen.width + screen.height + 
                   navigator.language + new Date().getTimezoneOffset();
        };

        // Self-defending mechanisms
        const originalCode = "%s";
        let executionAllowed = true;
        
        // Anti-tampering checks
        setInterval(() => {
            if(window.outerHeight - window.innerHeight > 100) {
                executionAllowed = false;
            }
        }, 1000);

        // Execution wrapper
        if(!checkVM() && executionAllowed) {
            return eval(originalCode);
        }
        return null;
    })();
    `
	return []byte(fmt.Sprintf(template, base64.StdEncoding.EncodeToString(code)))
}

// New method for string encryption
func (p *Pipeline) encryptStrings(code []byte) []byte {
	// Implementation for string encryption
	// This would identify and encrypt string literals in the code
	return code
}

// New method for control flow flattening
func (p *Pipeline) flattenControlFlow(code []byte) []byte {
	// Implementation for control flow flattening
	// This would restructure the code to make it harder to analyze
	return code
}

// New method for adding polymorphic code
func (p *Pipeline) addPolymorphicLayer(code []byte) []byte {
	template := `
        (function(){
            const _0x%x = function(code){
                return code.split('').reverse().join('');
            };
            const _0x%x = function(code){
                return code.split('').map(c => 
                    String.fromCharCode(c.charCodeAt(0) ^ 0x%x)
                ).join('');
            };
            return (function(){
                let result = %q;
                result = _0x%x(result);
                result = _0x%x(result);
                return result;
            })();
        })()
    `
	key1 := mrand.Int31() // Use mrand instead of rand
	key2 := mrand.Int31()
	key3 := mrand.Int31()
	return []byte(fmt.Sprintf(template,
		key1, key2, key3,
		string(code), key1, key2))
}

// New advanced loader
func (p *Pipeline) addAdvancedLoader(encrypted []byte) []byte {
	key := make([]byte, 32)
	crand.Read(key) // Use crand for cryptographic operations

	loader := fmt.Sprintf(`
    (function(){
        const _0x%x = %q;
        const _0x%x = %q;
        
        function decode(str) {
            // Complex decoding logic here
            return str;
        }
        
        function verify() {
            return (
                !window.Firebug && 
                !window.__REACT_DEVTOOLS_GLOBAL_HOOK__ &&
                !window.__REDUX_DEVTOOLS_EXTENSION__ &&
                window.chrome === undefined
            );
        }
        
        if(verify()) {
            const payload = decode(_0x%x);
            (new Function(payload))();
        }
    })();
    `,
		mrand.Int31(), key,
		mrand.Int31(), encrypted,
		mrand.Int31())

	return []byte(loader)
}

// Add new method for WASM wrapping
func (p *Pipeline) wrapWithWASM(code []byte) []byte {
	wasmLoader := `
        (async function(){
            try {
                const wasmInstance = await WebAssembly.instantiateStreaming(
                    fetch("/wasm/decoder.wasm"),
                    {
                        env: {
                            memory: new WebAssembly.Memory({ initial: 256 }),
                            abort: () => console.log("Abort called")
                        }
                    }
                );
                const decoder = wasmInstance.instance.exports.decode;
                return decoder(%q);
            } catch(e) {
                console.error("WASM execution failed");
                return null;
            }
        })()
    `
	return []byte(fmt.Sprintf(wasmLoader, code))
}

func init() {
	mrand.Seed(time.Now().UnixNano()) // Use mrand instead of rand
}
