(() => {
  // <stdin>
  document.addEventListener("DOMContentLoaded", () => {
    const couponButtons = document.querySelectorAll(".coupon-button");
    couponButtons.forEach((button) => {
      button.addEventListener("click", async (e) => {
        const code = button.getAttribute("data-code");
        const targetUrl = button.getAttribute("data-url");
        const originalText = button.querySelector("span").innerText;
        const textSpan = button.querySelector("span");
        try {
          await navigator.clipboard.writeText(code);
          textSpan.innerText = "Copied!";
          button.classList.add("bg-green-600", "text-white");
          button.classList.remove("bg-white/10");
          setTimeout(() => {
            window.location.href = targetUrl;
          }, 800);
        } catch (err) {
          console.error("Failed to copy text: ", err);
        }
      });
    });
  });
})();
