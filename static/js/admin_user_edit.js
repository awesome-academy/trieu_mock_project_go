document.addEventListener("DOMContentLoaded", function () {
  const updateUserBtn = document.getElementById("updateUserBtn");
  const editUserForm = document.getElementById("editUserForm");
  const userId = editUserForm.dataset.userId;

  // Initialize Skill Manager
  SkillManager.init({
    skillsListId: "skillsList",
    newSkillSelectId: "newSkillSelect",
    addSkillBtnId: "addSkillBtn",
  });

  // Update user
  updateUserBtn.addEventListener("click", async function () {
    const formData = new FormData(editUserForm);
    const data = {
      name: formData.get("name"),
      email: formData.get("email"),
      birthday: formData.get("birthday") || null,
      position_id: parseInt(formData.get("position_id")),
      team_id: formData.get("team_id")
        ? parseInt(formData.get("team_id"))
        : null,
      skills: SkillManager.getSelectedSkills(),
    };

    try {
      const response = await AdminUserService.updateUser(userId, data);
      Toast.success(response.message || "User updated successfully");
      setTimeout(() => {
        window.location.href = `/admin/users/${userId}`;
      }, 1500);
    } catch (error) {
      console.error("Error updating user:", error);
      let msg = error.message || "Failed to update user";
      if (error.details && typeof error.details === "object") {
        const details = Object.entries(error.details)
          .map(([field, err]) => `${field}: ${err}`)
          .join("<br>");
        msg += `<br><small>${details}</small>`;
      }
      Toast.error(msg);
    }
  });
});
