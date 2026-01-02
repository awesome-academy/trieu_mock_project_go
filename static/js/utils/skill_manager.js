/**
 * Skill Manager for User Create/Edit pages
 */
const SkillManager = {
  init: function (config) {
    this.skillsList = document.getElementById(config.skillsListId);
    this.newSkillSelect = document.getElementById(config.newSkillSelectId);
    this.addSkillBtn = document.getElementById(config.addSkillBtnId);

    if (!this.skillsList || !this.newSkillSelect || !this.addSkillBtn) return;

    this.addSkillBtn.addEventListener("click", () => this.handleAddSkill());
    this.skillsList.addEventListener("click", (e) => this.handleRemoveSkill(e));

    this.updateAvailableSkills();
  },

  updateAvailableSkills: function () {
    const currentSkillIds = Array.from(
      this.skillsList.querySelectorAll("tr")
    ).map((tr) => tr.dataset.skillId);

    Array.from(this.newSkillSelect.options).forEach((option) => {
      if (option.value === "") return;
      option.style.display = currentSkillIds.includes(option.value)
        ? "none"
        : "block";
    });
    this.newSkillSelect.value = "";
  },

  handleAddSkill: function () {
    const skillId = this.newSkillSelect.value;
    if (!skillId) return;

    const skillName =
      this.newSkillSelect.options[this.newSkillSelect.selectedIndex].dataset
        .name;

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

    this.skillsList.appendChild(tr);
    this.updateAvailableSkills();
  },

  handleRemoveSkill: function (e) {
    if (e.target.closest(".remove-skill-btn")) {
      e.target.closest("tr").remove();
      this.updateAvailableSkills();
    }
  },

  getSelectedSkills: function () {
    const skills = [];
    this.skillsList.querySelectorAll("tr").forEach((tr) => {
      skills.push({
        id: parseInt(tr.dataset.skillId),
        level: parseInt(tr.querySelector(".skill-level").value),
        used_year_number: parseInt(tr.querySelector(".skill-years").value),
      });
    });
    return skills;
  },
};
