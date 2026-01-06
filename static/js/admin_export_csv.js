document.addEventListener("DOMContentLoaded", function () {
  const exportForm = document.getElementById("exportForm");

  if (exportForm) {
    exportForm.addEventListener("submit", function (e) {
      const exportType = document.getElementById("exportType").value;
      if (!exportType) {
        e.preventDefault();
        Toast.error("Please select a data type to export");
        return;
      }

      // We don't prevent default here because we want the browser to handle the file download
      // But we can show a success message
      Toast.success("Exporting " + exportType + " data...");
    });
  }
});
