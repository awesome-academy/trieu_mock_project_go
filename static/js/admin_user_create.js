document.addEventListener("DOMContentLoaded", function () {
  const createUserBtn = document.getElementById("createUserBtn");
  const createUserForm = document.getElementById("createUserForm");

  // Initialize Skill Manager
  SkillManager.init({
    skillsListId: "skillsList",
    newSkillSelectId: "newSkillSelect",
    addSkillBtnId: "addSkillBtn",
  });

  // Create user
  createUserBtn.addEventListener("click", async function () {
    if (!createUserForm.checkValidity()) {
      createUserForm.reportValidity();
      return;
    }

    const formData = new FormData(createUserForm);
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
      const response = await AdminUserService.createUser(data);
      Toast.success(response.message || "User created successfully");
      setTimeout(() => {
        window.location.href = "/admin/users";
      }, 1500);
    } catch (error) {
      console.error("Error creating user:", error);
      let msg = error.message || "Failed to create user";
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
