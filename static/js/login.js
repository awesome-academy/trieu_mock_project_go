$(document).ready(function () {
  // Check if already logged in
  if (AuthService.isAuthenticated()) {
    // Already logged in, redirect to dashboard
    window.location.href = "/";
  }

  // Handle login form submission
  $("#loginForm").on("submit", async function (e) {
    e.preventDefault();

    const email = $("#email").val();
    const password = $("#password").val();
    const $loginBtn = $("#loginBtn");
    const $errorAlert = $("#errorAlert");
    const $errorMessage = $("#errorMessage");
    const $successAlert = $("#successAlert");

    // Hide alerts
    $errorAlert.addClass("d-none");
    $successAlert.addClass("d-none");

    // Disable button during request
    $loginBtn.prop("disabled", true);
    $loginBtn.html(
      '<span class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>Logging in...'
    );

    try {
      // Use AuthService for login
      await AuthService.login(email, password);

      // Show success message
      $successAlert.removeClass("d-none");

      // Redirect to dashboard after 1.5 seconds
      setTimeout(() => {
        window.location.href = "/";
      }, 1500);
    } catch (error) {
      // Show error message
      const errorMsg =
        error.responseJSON?.error ||
        error.message ||
        "An error occurred during login";
      $errorMessage.text(errorMsg);
      $errorAlert.removeClass("d-none");

      // Re-enable button
      $loginBtn.prop("disabled", false);
      $loginBtn.html("Login");
    }
  });
});
