// filepath: c:\Users\gagan\OneDrive\Desktop\decryptor\server\internal\security\antidebug.go
package security

import "fmt" // Add this import

type AntiDebug struct {
	enabled bool
}

func NewAntiDebug(enabled bool) *AntiDebug {
	return &AntiDebug{enabled: enabled}
}

func (a *AntiDebug) Protect(code []byte) []byte {
	if !a.enabled {
		return code
	}

	antiDebugWrapper := `
        (function(){
            const checkDebugger = function() {
                const start = performance.now();
                debugger;
                return performance.now() - start > 100;
            };
            
            if(checkDebugger()) {
                throw new Error("Debugging detected!");
            }
            
            // Original code here
            %s
        })();
    `

	return []byte(fmt.Sprintf(antiDebugWrapper, string(code))) // Changed sprintf to fmt.Sprintf
}
