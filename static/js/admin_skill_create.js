document.addEventListener("DOMContentLoaded", function () {
  const createSkillBtn = document.getElementById("createSkillBtn");
  const createSkillForm = document.getElementById("createSkillForm");
  if (createSkillBtn) {
    createSkillBtn.addEventListener("click", async function () {
      if (!createSkillForm.checkValidity()) {
        createSkillForm.reportValidity();
        return;
      }
      const formData = new FormData(createSkillForm);
      const data = {
        name: formData.get("name"),
      };
      try {
        const response = await AdminSkillService.createSkill(data);
        Toast.success(response.message || "Skill created successfully");
        setTimeout(() => {
          window.location.href = "/admin/skills";
        }, 1500);
      } catch (error) {
        console.error("Error creating skill:", error);
        Toast.error(error.message || "Failed to create skill");
      }
    });
  }
});
