document.addEventListener("DOMContentLoaded", function () {
  const skillsList = document.getElementById("skillsList");
  const newSkillSelect = document.getElementById("newSkillSelect");
  const addSkillBtn = document.getElementById("addSkillBtn");
  const updateUserBtn = document.getElementById("updateUserBtn");
  const editUserForm = document.getElementById("editUserForm");
  const userId = editUserForm.dataset.userId;

  // Function to update the available skills in the dropdown
  function updateAvailableSkills() {
    const currentSkillIds = Array.from(skillsList.querySelectorAll("tr")).map(
      (tr) => tr.dataset.skillId
    );

    Array.from(newSkillSelect.options).forEach((option) => {
      if (option.value === "") return;
      if (currentSkillIds.includes(option.value)) {
        option.style.display = "none";
      } else {
        option.style.display = "block";
      }
    });
    newSkillSelect.value = "";
  }

  // Initial update
  updateAvailableSkills();

  // Add skill
  addSkillBtn.addEventListener("click", function () {
    const skillId = newSkillSelect.value;
    if (!skillId) return;

    const skillName =
      newSkillSelect.options[newSkillSelect.selectedIndex].dataset.name;

    const tr = document.createElement("tr");
    tr.dataset.skillId = skillId;
    tr.innerHTML = `
      <td>${skillName}</td>
      <td>
        <select class="form-select form-select-sm skill-level">
          ${Array.from({ length: 10 }, (_, i) => i + 1)
            .map((i) => `<option value="${i}">${i}</option>`)
            .join("")}
        </select>
      </td>
      <td>
        <input type="number" class="form-control form-control-sm skill-years" 
          value="0" min="0" max="100">
      </td>
      <td>
        <button type="button" class="btn btn-outline-danger btn-sm remove-skill-btn">
          <i class="bi bi-trash"></i>
        </button>
      </td>
    `;

    skillsList.appendChild(tr);
    updateAvailableSkills();
  });

  // Remove skill
  skillsList.addEventListener("click", function (e) {
    if (e.target.closest(".remove-skill-btn")) {
      e.target.closest("tr").remove();
      updateAvailableSkills();
    }
  });

  // Update user
  updateUserBtn.addEventListener("click", async function () {
    const formData = new FormData(editUserForm);
    const data = {
      name: formData.get("name"),
      email: formData.get("email"),
      birthday: formData.get("birthday"),
      position_id: parseInt(formData.get("position_id")),
      team_id: formData.get("team_id")
        ? parseInt(formData.get("team_id"))
        : null,
      skills: [],
    };

    skillsList.querySelectorAll("tr").forEach((tr) => {
      const skillId = parseInt(tr.dataset.skillId);
      const level = parseInt(tr.querySelector(".skill-level").value);
      const years = parseInt(tr.querySelector(".skill-years").value);
      data.skills.push({
        id: skillId,
        level: level,
        used_year_number: years,
      });
    });

    try {
      const response = await fetch(`/admin/users/${userId}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      if (response.ok) {
        window.location.href = `/admin/users/${userId}`;
      } else {
        const err = await response.json();
        alert("Error: " + (err.error || "Failed to update user"));
      }
    } catch (error) {
      console.error("Error updating user:", error);
      alert("An error occurred while updating the user.");
    }
  });
});
