document.addEventListener("DOMContentLoaded", function () {
  const deleteUserBtn = document.getElementById("deleteUserBtn");

  if (deleteUserBtn) {
    deleteUserBtn.addEventListener("click", async function () {
      const userId = this.dataset.userId;
      const userName = this.dataset.userName;

      if (
        !confirm(
          `Are you sure you want to delete user "${userName}"? This action cannot be undone.`
        )
      ) {
        return;
      }

      try {
        const response = await AdminUserService.deleteUser(userId);
        Toast.success(response.message || "User deleted successfully");
        setTimeout(() => {
          window.location.href = "/admin/users";
        }, 1500);
      } catch (error) {
        console.error("Error deleting user:", error);
        Toast.error(error.message || "Failed to delete user");
      }
    });
  }
});
