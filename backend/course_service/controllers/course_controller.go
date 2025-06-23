package controllers

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/backlog-developer/video_content_management_system/backend/course_service/models"

	"github.com/gofiber/fiber/v2"
)

// GET /courses
func GetAllCourses(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, title, description, category_id, created_at, instructor_id FROM courses")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch courses"})
		}
		defer rows.Close()

		var courses []models.Course
		for rows.Next() {
			var course models.Course
			if err := rows.Scan(&course.ID, &course.Title, &course.Description, &course.CategoryID, &course.CreatedAt, &course.InstructorID); err != nil {
				fmt.Println("❌ Scan error in GetAllCourses:", err)
				return c.Status(500).JSON(fiber.Map{"error": "Failed to scan course data"})
			}
			courses = append(courses, course)
		}
		fmt.Println("✅ Courses fetched successfully")
		return c.JSON(courses)
	}
}

// GET /courses/:id
func GetCourseByID(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var course models.Course

		query := `SELECT id, title, description, category_id, created_at, instructor_id FROM courses WHERE id = $1`
		err := db.QueryRow(query, id).Scan(&course.ID, &course.Title, &course.Description, &course.CategoryID, &course.CreatedAt, &course.InstructorID)
		if err == sql.ErrNoRows {
			fmt.Println("⚠️ No course found with ID:", id)
			return c.Status(404).JSON(fiber.Map{"error": "Course not found"})
		} else if err != nil {
			fmt.Println("❌ DB error in GetCourseByID:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve course"})
		}

		fmt.Println("✅ Course retrieved:", course.ID)
		return c.JSON(course)
	}
}

// POST /courses
func CreateCourse(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var course models.Course
		if err := c.BodyParser(&course); err != nil {
			fmt.Println("❌ BodyParser error in CreateCourse:", err)
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		query := `
			INSERT INTO courses (title, description, category_id, created_at, instructor_id, created_by)
            VALUES ($1, $2, $3, NOW(), $4, $5) RETURNING id

		`
		err := db.QueryRow(query,
			course.Title,
			course.Description,
			course.CategoryID,
			course.InstructorID,
			course.CreatedBy,
		).Scan(&course.ID)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create course", "details": err.Error()})
		}

		fmt.Println("✅ Course created with ID:", course.ID)
		return c.Status(201).JSON(course)
	}
}

// PUT /courses/:id
func UpdateCourse(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var course models.Course
		if err := c.BodyParser(&course); err != nil {
			fmt.Println("❌ BodyParser error in UpdateCourse:", err)
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		courseID, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("❌ ID parse error in UpdateCourse:", err)
			return c.Status(400).JSON(fiber.Map{"error": "Invalid course ID"})
		}
		course.ID = courseID

		query := `
			UPDATE courses
			SET title = $1, description = $2, category_id = $3, instructor_id = $4
			WHERE id = $5
		`
		_, err = db.Exec(query, course.Title, course.Description, course.CategoryID, course.InstructorID, course.ID)
		if err != nil {
			fmt.Println("❌ Update error in UpdateCourse:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update course"})
		}

		fmt.Println("✅ Course updated:", course.ID)
		return c.JSON(fiber.Map{"message": "Course updated successfully"})
	}
}

// DELETE /courses/:id
func DeleteCourse(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		query := `DELETE FROM courses WHERE id = $1`
		result, err := db.Exec(query, id)
		if err != nil {
			fmt.Println("❌ Delete error in DeleteCourse:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete course"})
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			fmt.Println("⚠️ No course deleted. ID not found:", id)
			return c.Status(404).JSON(fiber.Map{"error": "Course not found"})
		}

		fmt.Println("✅ Course deleted with ID:", id)
		return c.JSON(fiber.Map{"message": "Course deleted successfully"})
	}
}
