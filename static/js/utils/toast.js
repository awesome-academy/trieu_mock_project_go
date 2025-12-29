/**
 * Toast Utility for showing notifications
 */
const Toast = {
  container: null,

  init: function () {
    if (this.container) return;

    this.container = document.createElement("div");
    this.container.className =
      "toast-container position-fixed top-0 start-50 translate-middle-x p-3";
    this.container.style.zIndex = "1055";
    document.body.appendChild(this.container);
  },

  show: function (message, type = "success") {
    this.init();

    const toastId = "toast-" + Date.now();
    const bgClass = type === "success" ? "bg-success" : "bg-danger";
    const iconClass =
      type === "success"
        ? "bi-check-circle-fill"
        : "bi-exclamation-triangle-fill";

    const toastHTML = `
      <div id="${toastId}" class="toast align-items-center text-white ${bgClass} border-0" role="alert" aria-live="assertive" aria-atomic="true">
        <div class="d-flex">
          <div class="toast-body">
            <i class="bi ${iconClass} me-2"></i>
            ${message}
          </div>
          <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>
        </div>
      </div>
    `;

    this.container.insertAdjacentHTML("beforeend", toastHTML);
    const toastElement = document.getElementById(toastId);
    const bsToast = new bootstrap.Toast(toastElement, { delay: 5000 });
    bsToast.show();

    toastElement.addEventListener("hidden.bs.toast", () => {
      toastElement.remove();
    });
  },

  success: function (message) {
    this.show(message, "success");
  },

  error: function (message) {
    this.show(message, "error");
  },
};
