/**
 * Admin API Utility for handling session-based AJAX requests with CSRF
 */
const AdminAPI = {
  /**
   * Base request handler
   * @param {Object} options - jQuery AJAX options
   * @returns {Promise}
   */
  request: function (options) {
    const csrfToken = $('meta[name="csrf-token"]').attr("content");

    const defaultOptions = {
      contentType: "application/json",
      dataType: "json",
      headers: {},
    };

    if (csrfToken) {
      defaultOptions.headers["X-CSRF-TOKEN"] = csrfToken;
    }

    const ajaxOptions = $.extend(true, {}, defaultOptions, options);

    return $.ajax(ajaxOptions).catch((xhr) => {
      if (xhr.status === 401) {
        // Unauthorized for admin - redirect to admin login
        window.location.href = "/admin/login";
      }
      // Reject with the JSON response if available, otherwise the XHR object
      throw xhr.responseJSON || xhr;
    });
  },

  get: function (url, options = {}) {
    return this.request({ ...options, url, method: "GET" });
  },

  post: function (url, data, options = {}) {
    return this.request({
      ...options,
      url,
      method: "POST",
      data: JSON.stringify(data),
    });
  },

  put: function (url, data, options = {}) {
    return this.request({
      ...options,
      url,
      method: "PUT",
      data: JSON.stringify(data),
    });
  },

  delete: function (url, options = {}) {
    return this.request({ ...options, url, method: "DELETE" });
  },
};
