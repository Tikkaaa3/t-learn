import { apiClient } from "./client";

export interface Course {
  id: string;
  title: string;
  description: string;
}

export interface Lesson {
  id: string;
  title: string;
  position: number;
  completed: boolean;
}

export interface TaskStep {
  position: number;
  command: string;
  expected_output: string;
}

export interface TaskResponse {
  lesson_id: string;
  lesson_title: string;
  lesson_content: string;
  task_id: string;
  task_description: string;
  steps: TaskStep[];
}

export async function getCourses(): Promise<Course[]> {
  return apiClient<Course[]>("/courses");
}

export async function getLessons(courseId: string): Promise<Lesson[]> {
  return apiClient<Lesson[]>(`/courses/${courseId}/lessons`);
}

export async function getTask(lessonId: string): Promise<TaskResponse> {
  return apiClient<TaskResponse>(`/lessons/${lessonId}/task`);
}

// --- ADMIN API ---

export async function createCourse(title: string, description: string) {
  return apiClient("/admin/courses", {
    method: "POST",
    body: JSON.stringify({ title, description }),
  });
}

export async function deleteCourse(courseId: string) {
  return apiClient(`/admin/courses/${courseId}`, {
    method: "DELETE",
  });
}

export async function createLesson(
  courseId: string,
  title: string,
  content: string,
  position: number,
) {
  return apiClient(`/admin/courses/${courseId}/lessons`, {
    method: "POST",
    body: JSON.stringify({ title, content, position }),
  });
}

export async function deleteLesson(lessonId: string) {
  return apiClient(`/admin/lessons/${lessonId}`, {
    method: "DELETE",
  });
}
