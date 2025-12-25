/**
 * API Utility for handling AJAX requests
 */
const API = {
  /**
   * Base request handler
   * @param {Object} options - jQuery AJAX options
   * @returns {Promise}
   */
  request: function (options) {
    const token = localStorage.getItem("accessToken");

    const defaultOptions = {
      contentType: "application/json",
      dataType: "json",
      headers: {},
    };

    if (token) {
      defaultOptions.headers["Authorization"] = `Bearer ${token}`;
    }

    const ajaxOptions = $.extend(true, {}, defaultOptions, options);

    return $.ajax(ajaxOptions).catch((xhr) => {
      if (xhr.status === 401) {
        // Unauthorized - clear token and redirect to login
        localStorage.removeItem("accessToken");
        localStorage.removeItem("userEmail");
        localStorage.removeItem("userName");
        localStorage.removeItem("userId");
        window.location.href = "/login";
      }
      throw xhr;
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
