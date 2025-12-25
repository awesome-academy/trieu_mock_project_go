$(document).ready(function () {
  // Check if already logged in
  const token = localStorage.getItem("accessToken");
  if (token) {
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
      // Prepare request body
      const requestBody = {
        user: {
          email: email,
          password: password,
        },
      };

      // Send login request
      const response = await $.ajax({
        type: "POST",
        url: "/login",
        data: JSON.stringify(requestBody),
        contentType: "application/json",
        dataType: "json",
      });

      // Extract access token from response
      const accessToken = response.user.access_token;

      if (!accessToken) {
        throw new Error("No access token received from server");
      }

      // Store token in localStorage
      localStorage.setItem("accessToken", accessToken);
      localStorage.setItem("userEmail", response.user.email);
      localStorage.setItem("userName", response.user.name);
      localStorage.setItem("userId", response.user.id);

      // Show success message
      $successAlert.removeClass("d-none");

      // Redirect to dashboard after 1.5 seconds
      window.location.href = "/";
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
