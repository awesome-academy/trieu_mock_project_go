document.addEventListener("DOMContentLoaded", function () {
  const updateSkillBtn = document.getElementById("updateSkillBtn");
  const editSkillForm = document.getElementById("editSkillForm");

  if (!updateSkillBtn) return;

  updateSkillBtn.addEventListener("click", async function () {
    if (!editSkillForm.checkValidity()) {
      editSkillForm.reportValidity();
      return;
    }

    const skillId = editSkillForm.getAttribute("data-id");
    const formData = new FormData(editSkillForm);
    const data = {
      name: formData.get("name"),
    };

    try {
      const response = await AdminSkillService.updateSkill(skillId, data);
      Toast.success(response.message || "Skill updated successfully");
      setTimeout(() => {
        window.location.href = "/admin/skills";
      }, 1500);
    } catch (error) {
      console.error("Error updating skill:", error);
      Toast.error(error.message || "Failed to update skill");
    }
  });
});
