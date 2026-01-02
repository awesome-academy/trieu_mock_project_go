document.addEventListener("DOMContentLoaded", function () {
  const updatePositionBtn = document.getElementById("updatePositionBtn");
  const editPositionForm = document.getElementById("editPositionForm");

  if (!updatePositionBtn) return;

  updatePositionBtn.addEventListener("click", async function () {
    if (!editPositionForm.checkValidity()) {
      editPositionForm.reportValidity();
      return;
    }

    const positionId = editPositionForm.getAttribute("data-id");
    const formData = new FormData(editPositionForm);
    const data = {
      name: formData.get("name"),
      abbreviation: formData.get("abbreviation"),
    };

    try {
      const response = await AdminPositionService.updatePosition(
        positionId,
        data
      );
      Toast.success(response.message || "Position updated successfully");
      setTimeout(() => {
        window.location.href = "/admin/positions";
      }, 1500);
    } catch (error) {
      console.error("Error updating position:", error);
      Toast.error(error.message || "Failed to update position");
    }
  });
});
