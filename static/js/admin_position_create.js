document.addEventListener("DOMContentLoaded", function () {
  const createPositionBtn = document.getElementById("createPositionBtn");
  const createPositionForm = document.getElementById("createPositionForm");

  createPositionBtn.addEventListener("click", async function () {
    if (!createPositionForm.checkValidity()) {
      createPositionForm.reportValidity();
      return;
    }

    const formData = new FormData(createPositionForm);
    const data = {
      name: formData.get("name"),
      abbreviation: formData.get("abbreviation"),
    };

    try {
      const response = await AdminPositionService.createPosition(data);
      Toast.success(response.message || "Position created successfully");
      setTimeout(() => {
        window.location.href = "/admin/positions";
      }, 1500);
    } catch (error) {
      console.error("Error creating position:", error);
      Toast.error(error.message || "Failed to create position");
    }
  });
});
