package database

import (
	"context"
	"fmt"
	"testing"

	"github.com/keegancsmith/sqlf"
	"github.com/sourcegraph/log/logtest"
	"github.com/stretchr/testify/require"

	"github.com/sourcegraph/sourcegraph/internal/database/dbtest"
	"github.com/sourcegraph/sourcegraph/internal/types"
)

func TestUserRoleAssign(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := logtest.Scoped(t)
	db := NewDB(logger, dbtest.NewDB(logger, t))
	store := db.UserRoles()

	user, role := createUserAndRole(ctx, t, db)

	t.Run("without user id", func(t *testing.T) {
		err := store.Assign(ctx, AssignUserRoleOpts{
			RoleID: role.ID,
		})
		require.Error(t, err)
		require.Equal(t, err.Error(), "missing user id")
	})

	t.Run("without role id", func(t *testing.T) {
		err := store.Assign(ctx, AssignUserRoleOpts{
			UserID: user.ID,
		})
		require.Error(t, err)
		require.Equal(t, err.Error(), "missing role id")
	})

	t.Run("success", func(t *testing.T) {
		err := store.Assign(ctx, AssignUserRoleOpts{
			RoleID: role.ID,
			UserID: user.ID,
		})
		require.NoError(t, err)

		ur, err := store.GetByRoleIDAndUserID(ctx, GetUserRoleOpts{
			RoleID: role.ID,
			UserID: user.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, ur)
		require.Equal(t, ur.RoleID, role.ID)
		require.Equal(t, ur.UserID, user.ID)

		// shoudln't fail the second time, since we are "upsert"-ing here
		err = store.Assign(ctx, AssignUserRoleOpts{
			RoleID: role.ID,
			UserID: user.ID,
		})
		require.NoError(t, err)
	})
}

func TestUserRoleBulkAssignForUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := logtest.Scoped(t)
	db := NewDB(logger, dbtest.NewDB(logger, t))
	store := db.UserRoles()

	user, role := createUserAndRole(ctx, t, db)
	role2, err := createTestRole(ctx, "another-role", false, t, db.Roles())
	require.NoError(t, err)

	t.Run("without user id", func(t *testing.T) {
		err := store.BulkAssignToUser(ctx, BulkAssignToUserOpts{})

		require.Error(t, err)
		require.Equal(t, err.Error(), "missing user id")
	})

	t.Run("without role ids", func(t *testing.T) {
		err := store.BulkAssignToUser(ctx, BulkAssignToUserOpts{
			UserID: user.ID,
		})

		require.Error(t, err)
		require.Equal(t, err.Error(), "missing role ids")
	})

	t.Run("success", func(t *testing.T) {
		roleIDs := []int32{role.ID, role2.ID}
		err := store.BulkAssignToUser(ctx, BulkAssignToUserOpts{
			UserID:  user.ID,
			RoleIDs: roleIDs,
		})
		require.NoError(t, err)

		urs, err := store.GetByUserID(ctx, GetUserRoleOpts{
			UserID: user.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, urs)
		require.Len(t, urs, len(roleIDs))
		for i, ur := range urs {
			require.Equal(t, ur.UserID, user.ID)
			require.Equal(t, ur.RoleID, roleIDs[i])
		}

		// shoudln't fail the second time, since we are "upsert"-ing here
		err = store.BulkAssignToUser(ctx, BulkAssignToUserOpts{
			UserID:  user.ID,
			RoleIDs: roleIDs,
		})
		require.NoError(t, err)
	})
}

func TestUserRoleAssignSysemRole(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := logtest.Scoped(t)
	db := NewDB(logger, dbtest.NewDB(logger, t))
	store := db.UserRoles()

	user, _ := createUserAndRole(ctx, t, db)

	t.Run("without user id", func(t *testing.T) {
		err := store.AssignSystemRole(ctx, AssignSystemRoleOpts{})
		require.ErrorContains(t, err, "user id is required")
	})

	t.Run("without role", func(t *testing.T) {
		err := store.AssignSystemRole(ctx, AssignSystemRoleOpts{
			UserID: user.ID,
		})
		require.ErrorContains(t, err, "role is required")
	})

	t.Run("success", func(t *testing.T) {
		err := store.AssignSystemRole(ctx, AssignSystemRoleOpts{
			UserID: user.ID,
			Role:   types.UserSystemRole,
		})
		require.NoError(t, err)

		urs, err := store.GetByUserID(ctx, GetUserRoleOpts{
			UserID: user.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, urs)
		require.Len(t, urs, 1)

		// shoudln't fail the second time, since we are "upsert"-ing here
		err = store.AssignSystemRole(ctx, AssignSystemRoleOpts{
			UserID: user.ID,
			Role:   types.UserSystemRole,
		})
		require.NoError(t, err)
	})
}

func TestUserRoleBulkAssignSystemRolesToUsers(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := logtest.Scoped(t)
	db := NewDB(logger, dbtest.NewDB(logger, t))
	store := db.UserRoles()

	user, _ := createUserAndRole(ctx, t, db)

	t.Run("without user id", func(t *testing.T) {
		err := store.BulkAssignSystemRolesToUser(ctx, BulkAssignSystemRolesToUserOpts{})
		require.ErrorContains(t, err, "user id is required")
	})

	t.Run("without roles", func(t *testing.T) {
		err := store.BulkAssignSystemRolesToUser(ctx, BulkAssignSystemRolesToUserOpts{
			UserID: user.ID,
		})
		require.ErrorContains(t, err, "roles are required")
	})

	t.Run("success", func(t *testing.T) {
		systemRoles := []types.SystemRole{types.SiteAdministratorSystemRole, types.UserSystemRole}
		err := store.BulkAssignSystemRolesToUser(ctx, BulkAssignSystemRolesToUserOpts{
			UserID: user.ID,
			Roles:  systemRoles,
		})
		require.NoError(t, err)

		urs, err := store.GetByUserID(ctx, GetUserRoleOpts{
			UserID: user.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, urs)
		require.Len(t, urs, len(systemRoles))

		// This shoudln't fail the second time since we are upserting.
		err = store.BulkAssignSystemRolesToUser(ctx, BulkAssignSystemRolesToUserOpts{
			UserID: user.ID,
			Roles:  systemRoles,
		})
		require.NoError(t, err)
	})
}

func TestUserRoleRevoke(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := logtest.Scoped(t)
	db := NewDB(logger, dbtest.NewDB(logger, t))
	store := db.UserRoles()

	user, role := createUserAndRole(ctx, t, db)

	// create a user role
	err := store.Assign(ctx, AssignUserRoleOpts{
		RoleID: role.ID,
		UserID: user.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("missing user id", func(t *testing.T) {
		err := store.Revoke(ctx, RevokeUserRoleOpts{
			RoleID: role.ID,
		})
		require.ErrorContains(t, err, "missing user id")
	})

	t.Run("missing role id", func(t *testing.T) {
		err := store.Revoke(ctx, RevokeUserRoleOpts{
			UserID: user.ID,
		})
		require.ErrorContains(t, err, "missing role id")
	})

	t.Run("with existing user role", func(t *testing.T) {
		err := store.Revoke(ctx, RevokeUserRoleOpts{
			RoleID: role.ID,
			UserID: user.ID,
		})
		require.NoError(t, err)

		ur, err := store.GetByRoleIDAndUserID(ctx, GetUserRoleOpts{
			RoleID: role.ID,
			UserID: user.ID,
		})
		require.Nil(t, ur)
		require.Error(t, err)
		require.Equal(t, err, &UserRoleNotFoundErr{
			RoleID: role.ID,
			UserID: user.ID,
		})
	})

	t.Run("with non-existent user role", func(t *testing.T) {
		roleID := int32(1234)
		userID := int32(4321)

		err := store.Revoke(ctx, RevokeUserRoleOpts{
			RoleID: roleID,
			UserID: userID,
		})
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to revoke user role")
	})
}

func TestUserRoleGetByRoleID(t *testing.T) {
	ctx := context.Background()
	logger := logtest.Scoped(t)
	db := NewDB(logger, dbtest.NewDB(logger, t))
	store := db.UserRoles()

	role := createTestRoleForUserRole(ctx, "TESTROLE", t, db)

	totalUsersWithRole := 10
	for i := 1; i <= totalUsersWithRole; i++ {
		username := fmt.Sprintf("ANOTHERTESTUSER%d", i)
		user := createTestUserWithoutRoles(t, db, username, false)

		err := store.Assign(ctx, AssignUserRoleOpts{
			RoleID: role.ID,
			UserID: user.ID,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Run("missing role id", func(t *testing.T) {
		urs, err := store.GetByRoleID(ctx, GetUserRoleOpts{})
		require.Error(t, err)
		require.Nil(t, urs)
		require.Equal(t, err.Error(), "missing role id")
	})

	t.Run("with provided role id", func(t *testing.T) {
		urs, err := store.GetByRoleID(ctx, GetUserRoleOpts{
			RoleID: role.ID,
		})

		require.NoError(t, err)
		require.Len(t, urs, totalUsersWithRole)

		for _, ur := range urs {
			require.Equal(t, ur.RoleID, role.ID)
		}
	})
}

func TestUserRoleGetByUserID(t *testing.T) {
	ctx := context.Background()
	logger := logtest.Scoped(t)
	db := NewDB(logger, dbtest.NewDB(logger, t))
	store := db.UserRoles()

	user := createTestUserWithoutRoles(t, db, "ANOTHERTESTUSER", false)

	totalRoles := 3
	for i := 1; i <= totalRoles; i++ {
		name := fmt.Sprintf("TESTROLE%d", i)
		role := createTestRoleForUserRole(ctx, name, t, db)

		err := store.Assign(ctx, AssignUserRoleOpts{
			RoleID: role.ID,
			UserID: user.ID,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Run("missing user id", func(t *testing.T) {
		urs, err := store.GetByUserID(ctx, GetUserRoleOpts{})
		require.Error(t, err)
		require.Nil(t, urs)
		require.Equal(t, err.Error(), "missing user id")
	})

	t.Run("with provided role id", func(t *testing.T) {
		urs, err := store.GetByUserID(ctx, GetUserRoleOpts{
			UserID: user.ID,
		})

		require.NoError(t, err)
		require.Len(t, urs, totalRoles)

		for _, ur := range urs {
			require.Equal(t, ur.UserID, user.ID)
		}
	})
}

func TestUserRoleGetByRoleIDAndUserID(t *testing.T) {
	ctx := context.Background()
	logger := logtest.Scoped(t)
	db := NewDB(logger, dbtest.NewDB(logger, t))
	store := db.UserRoles()

	user, role := createUserAndRole(ctx, t, db)
	err := store.Assign(ctx, AssignUserRoleOpts{
		RoleID: role.ID,
		UserID: user.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("without user id", func(t *testing.T) {
		ur, err := store.GetByRoleIDAndUserID(ctx, GetUserRoleOpts{
			RoleID: role.ID,
		})
		require.Nil(t, ur)
		require.Error(t, err)
		require.Equal(t, err.Error(), "missing user id")
	})

	t.Run("without role id", func(t *testing.T) {
		ur, err := store.GetByRoleIDAndUserID(ctx, GetUserRoleOpts{
			UserID: user.ID,
		})
		require.Nil(t, ur)
		require.Error(t, err)
		require.Equal(t, err.Error(), "missing role id")
	})

	t.Run("with correct args", func(t *testing.T) {
		ur, err := store.GetByRoleIDAndUserID(ctx, GetUserRoleOpts{
			UserID: user.ID,
			RoleID: role.ID,
		})

		require.NoError(t, err)
		require.Equal(t, ur.RoleID, role.ID)
		require.Equal(t, ur.UserID, user.ID)
	})
}

func createUserAndRole(ctx context.Context, t *testing.T, db DB) (*types.User, *types.Role) {
	t.Helper()
	user := createTestUserWithoutRoles(t, db, "u1", false)
	role := createTestRoleForUserRole(ctx, "ANOTHERTESTROLE - 1", t, db)
	return user, role
}

func createTestRoleForUserRole(ctx context.Context, name string, t *testing.T, db DB) *types.Role {
	t.Helper()
	role, err := db.Roles().Create(ctx, name, false)
	if err != nil {
		t.Fatal(err)
	}
	return role
}

func createTestUserWithoutRoles(t *testing.T, db DB, username string, siteAdmin bool) *types.User {
	t.Helper()

	user := &types.User{
		Username:    username,
		DisplayName: "testuser",
	}

	q := sqlf.Sprintf("INSERT INTO users (username, site_admin) VALUES (%s, %t) RETURNING id, site_admin", user.Username, siteAdmin)
	err := db.QueryRowContext(context.Background(), q.Query(sqlf.PostgresBindVar), q.Args()...).Scan(&user.ID, &user.SiteAdmin)
	if err != nil {
		t.Fatal(err)
	}

	if user.SiteAdmin != siteAdmin {
		t.Fatalf("user.SiteAdmin=%t, but expected is %t", user.SiteAdmin, siteAdmin)
	}

	_, err = db.ExecContext(context.Background(), "INSERT INTO names(name, user_id) VALUES($1, $2)", user.Username, user.ID)
	if err != nil {
		t.Fatalf("failed to create name: %s", err)
	}

	return user
}
