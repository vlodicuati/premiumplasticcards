document.addEventListener('DOMContentLoaded', () => {
    // 1. Select all coupon buttons
    const couponButtons = document.querySelectorAll('.coupon-button');

    // 2. Add click event listener to each button
    couponButtons.forEach(button => {
        button.addEventListener('click', async (e) => {
            const code = button.getAttribute('data-code');
            const targetUrl = button.getAttribute('data-url');
            const originalText = button.querySelector('span').innerText;
            const textSpan = button.querySelector('span');

            try {
                // 3. Write code to clipboard
                await navigator.clipboard.writeText(code);
                
                // 4. Update UI to show success
                textSpan.innerText = 'Copied!';
                button.classList.add('bg-green-600', 'text-white');
                button.classList.remove('bg-white/10');

                // 5. Redirect to Luxor after brief delay
                setTimeout(() => {
                    window.location.href = targetUrl;
                }, 800);

            } catch (err) {
                console.error('Failed to copy text: ', err);
                // Fallback for older browsers could go here
            }
        });
    });
});