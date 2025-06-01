async function loadProtectedCode() {
    try {
        const response = await fetch('http://localhost:8080/js/protected.js');
        const code = await response.text();
        
        // Create script element
        const script = document.createElement('script');
        script.text = code;
        document.body.appendChild(script);
        
        return true;
    } catch (error) {
        console.error('Failed to load protected code:', error);
        return false;
    }
}