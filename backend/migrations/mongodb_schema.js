// MongoDB Schema for users collection
db.createCollection("users", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["username", "nickname", "password", "role"],
      properties: {
        username: {
          bsonType: "string",
          minLength: 3,
          maxLength: 50,
          description: "Username must be unique and 3-50 characters"
        },
        nickname: {
          bsonType: "string",
          minLength: 1,
          maxLength: 50,
          description: "Nickname must be unique and 1-50 characters"
        },
        password: {
          bsonType: "string",
          description: "Bcrypt hashed password"
        },
        role: {
          bsonType: "string",
          enum: ["system_admin", "project_manager", "project_member"],
          description: "User role"
        },
        createdAt: {
          bsonType: "date",
          description: "Creation timestamp"
        },
        updatedAt: {
          bsonType: "date",
          description: "Last update timestamp"
        },
        deletedAt: {
          bsonType: ["date", "null"],
          description: "Soft delete timestamp"
        }
      }
    }
  }
});

// Create unique indexes
db.users.createIndex({ username: 1 }, { unique: true });
db.users.createIndex({ nickname: 1 }, { unique: true });
db.users.createIndex({ role: 1 });
db.users.createIndex({ deletedAt: 1 });
