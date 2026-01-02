const AdminTeamService = {
  listTeams: async function (params) {
    const query = new URLSearchParams(params).toString();
    const response = await fetch(`/admin/teams/partial/search?${query}`);
    if (!response.ok) throw new Error("Failed to fetch teams");
    return await response.text();
  },

  createTeam: async function (data) {
    const response = await fetch("/admin/teams", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-CSRF-Token": document.querySelector('meta[name="csrf-token"]')
          .content,
      },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      try {
        const error = await response.json();
        throw new Error(error.message || "Failed to create team");
      } catch {
        throw new Error("Failed to create team");
      }
    }
    return await response.json();
  },

  updateTeam: async function (id, data) {
    const response = await fetch(`/admin/teams/${id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "X-CSRF-Token": document.querySelector('meta[name="csrf-token"]')
          .content,
      },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || "Failed to update team");
    }
    return await response.json();
  },

  deleteTeam: async function (id) {
    const response = await fetch(`/admin/teams/${id}`, {
      method: "DELETE",
      headers: {
        "X-CSRF-Token": document.querySelector('meta[name="csrf-token"]')
          .content,
      },
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || "Failed to delete team");
    }
    return await response.json();
  },

  addMember: async function (teamId, userId) {
    const response = await fetch(`/admin/teams/${teamId}/members`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-CSRF-Token": document.querySelector('meta[name="csrf-token"]')
          .content,
      },
      body: JSON.stringify({ user_id: parseInt(userId) }),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || "Failed to add member");
    }
    return await response.json();
  },

  removeMember: async function (teamId, userId) {
    const response = await fetch(`/admin/teams/${teamId}/members/${userId}`, {
      method: "DELETE",
      headers: {
        "X-CSRF-Token": document.querySelector('meta[name="csrf-token"]')
          .content,
      },
    });
    if (!response.ok) {
      try {
        const error = await response.json();
        throw new Error(error.message || "Failed to remove member");
      } catch {
        throw new Error("Failed to remove member");
      }
    }
    return await response.json();
  },
};
