import type { CommandDefinition } from "../types";
import type { Course, Lesson } from "../api/content";
import { loginUser, generateApiKey, registerUser } from "../api/auth";
import {
  getCourses,
  getLessons,
  getTask,
  createCourse,
  deleteCourse,
  createLesson,
  deleteLesson,
} from "../api/content";

// --- GLOBAL STATE ---
interface ShellState {
  user: string | null;
  path: string[]; // e.g. ["Go Mastery"]
  cachedCourses: Course[];
  cachedLessons: Lesson[];
}

// Initial State
const state: ShellState = {
  user: localStorage.getItem("t_learn_user") || null,
  path: [],
  cachedCourses: [],
  cachedLessons: [],
};

// --- HELPER: Resolve ID by Exact ID or Fuzzy Name ---
function resolveId(
  query: string,
  list: { id: string; title: string }[],
): string | null {
  // Exact ID match?
  if (list.find((item) => item.id === query)) return query;

  // Fuzzy Name match?
  const lowerQuery = query.toLowerCase();
  const found = list.find((item) =>
    item.title.toLowerCase().includes(lowerQuery),
  );
  return found ? found.id : null;
}

// --- EXPORT FOR REACT UI ---
export const getPrompt = () => {
  if (state.path.length > 0) {
    const currentDir = state.path[state.path.length - 1];
    return `${currentDir} $`;
  }
  return "$";
};

// --- COMMAND DEFINITIONS ---

const help: CommandDefinition = {
  description: "Show available commands",
  execute: async () => {
    return {
      type: "info",
      output: `
### Available Commands
\`\`\`text
  help                          - Show this message
  clear                         - Clear the terminal
  register <user> <mail> <pass> - Create account
  login <user> <pass>           - Log in
  logout                        - Log out
  whoami                        - Show current user
  token                         - Generate CLI API Key
  courses                       - List available courses
  lessons <course_name>         - Enter a course
  start <lesson_name>           - Start a lesson task
\`\`\`
`,
    };
  },
};

const clear: CommandDefinition = {
  description: "Clear terminal",
  execute: async () => {
    return { type: "clear", output: "" };
  },
};

const register: CommandDefinition = {
  description: "Create a new account",
  execute: async (args) => {
    if (args.length < 3)
      return {
        type: "error",
        output: "Usage: register <username> <email> <password>",
      };
    const [username, email, password] = args;
    try {
      await registerUser(username, email, password);
      return {
        type: "success",
        output: `Account created for ${username}!\nYou can now log in.`,
      };
    } catch (err: any) {
      return { type: "error", output: `Registration failed: ${err.message}` };
    }
  },
};

const login: CommandDefinition = {
  description: "Log in",
  execute: async (args) => {
    if (args.length < 2)
      return { type: "error", output: "Usage: login <user> <pass>" };
    const [username, password] = args;
    try {
      const data = await loginUser(username, password);
      localStorage.setItem("t_learn_token", data.token);

      // Update State
      state.user = username;
      localStorage.setItem("t_learn_user", username);

      return { type: "success", output: `Logged in as ${username}.` };
    } catch (err: any) {
      return { type: "error", output: `Login failed: ${err.message}` };
    }
  },
};

const logout: CommandDefinition = {
  description: "Log out of the session",
  execute: async () => {
    // Clear the storage
    localStorage.removeItem("t_learn_token");
    localStorage.removeItem("t_learn_user");

    // Reset the internal state
    state.user = null;
    state.path = []; // Optional: Reset path to root

    return { type: "success", output: "Logged out successfully." };
  },
};

const whoami: CommandDefinition = {
  description: "Show current user",
  execute: async () => {
    return state.user
      ? { type: "success", output: state.user }
      : { type: "error", output: "Not logged in." };
  },
};

const token: CommandDefinition = {
  description: "Generate CLI API Key",
  execute: async () => {
    try {
      const data = await generateApiKey();

      return {
        type: "success",
        output:
          `### 🔑 API Key Generated\n` +
          `Use this token to authenticate your CLI tool:\n\n` +
          `\`\`\`bash\nt-cli login ${data.api_key}\n\`\`\`\n` +
          `_Keep this token safe!_`,
      };
    } catch (err: any) {
      return { type: "error", output: `Failed: ${err.message}` };
    }
  },
};

const courses: CommandDefinition = {
  description: "List available courses",
  execute: async () => {
    try {
      const courses = await getCourses();
      state.cachedCourses = courses;

      if (courses.length === 0)
        return { type: "info", output: "No courses found." };

      const list = courses
        .map((c) => `- **${c.title}**`) // Bold the title
        .join("\n");

      return { type: "info", output: `### Available Courses:\n${list}` };
    } catch (err: any) {
      return {
        type: "error",
        output: `Failed to fetch courses: ${err.message}`,
      };
    }
  },
};

const lessons: CommandDefinition = {
  description: "List lessons in a course",
  execute: async (args) => {
    if (args.length < 1)
      return { type: "error", output: "Usage: lessons <course_name>" };

    const courseQuery = args.join(" "); // Handle "Rust Basics"

    // Ensure cache
    if (state.cachedCourses.length === 0) {
      try {
        state.cachedCourses = await getCourses();
      } catch (e) {}
    }

    const courseId = resolveId(courseQuery, state.cachedCourses);
    if (!courseId)
      return { type: "error", output: `Course '${courseQuery}' not found.` };

    try {
      const lessons = await getLessons(courseId);
      state.cachedLessons = lessons;

      if (lessons.length === 0)
        return { type: "info", output: `No lessons in '${courseQuery}'.` };

      const list = lessons
        .map((l) => {
          const mark = l.completed ? "x" : " "; // x for done, space for todo
          return `- [${mark}] ${l.title}`;
        })
        .join("\n");

      return {
        type: "info",
        output: `### Lessons in ${courseQuery}:\n${list}`,
      };
    } catch (err: any) {
      return { type: "error", output: `Failed: ${err.message}` };
    }
  },
};

const start: CommandDefinition = {
  description: "Start a lesson task",
  execute: async (args) => {
    if (args.length < 1)
      return { type: "error", output: "Usage: start <lesson_name>" };

    const query = args.join(" ");

    // Resolve Lesson ID
    const lessonId = resolveId(query, state.cachedLessons);
    if (!lessonId) {
      return {
        type: "error",
        output: `Lesson '${query}' not found.\n(Did you run 'lessons <course>' first?)`,
      };
    }

    try {
      // Use your EXISTING getTask function
      const data = await getTask(lessonId);

      // Build Rich Markdown Output
      let output = `# ${data.lesson_title}\n\n`;

      // The content (includes the CLI Helper we added in the seeder)
      output += `${data.lesson_content}\n\n`;

      // --- Task Section ---
      output += `## 🎯 Your Task\n`;
      output += `${data.task_description}\n\n`;

      if (data.steps && data.steps.length > 0) {
        output += `**Steps to execute:**\n`;
        data.steps.forEach((step: any) => {
          // Render commands as inline code blocks
          output += `${step.position}. \`${step.command}\`\n`;
        });
      }

      // --- Verification Section ---
      output += `\n---\n`;
      output += `### ✅ Verification\n`;
      output += `Run this command to check your work:\n`;

      // The Copy-Paste Block
      output += `\`\`\`bash\nt-cli ${data.lesson_id}\n\`\`\``;

      return { type: "info", output };
    } catch (err: any) {
      return { type: "error", output: `Failed to load task: ${err.message}` };
    }
  },
};

// --- ADMIN COMMANDS (SMART VERSIONS) ---

const mkcourse: CommandDefinition = {
  description: "Create a course (Admin)",
  execute: async (args) => {
    // Usage: mkcourse "Go Mastery" "Learn Go"
    if (args.length < 2)
      return {
        type: "error",
        output: 'Usage: mkcourse "<Title>" "<Description>"',
      };

    const [title, desc] = args;
    try {
      await createCourse(title, desc);
      state.cachedCourses = await getCourses(); // Refresh cache immediately
      return { type: "success", output: `Course "${title}" created.` };
    } catch (err: any) {
      return { type: "error", output: `Failed: ${err.message}` };
    }
  },
};

const rmcourse: CommandDefinition = {
  description: "Delete a course by Name or ID (Admin)",
  execute: async (args) => {
    if (args.length < 1)
      return { type: "error", output: "Usage: rmcourse <course_name_or_id>" };

    const query = args.join(" "); // Handle names with spaces like "Go Mastery"

    // Ensure we have the list to look up names
    if (state.cachedCourses.length === 0) {
      try {
        state.cachedCourses = await getCourses();
      } catch (e) {}
    }

    // Resolve Name -> ID
    const courseId = resolveId(query, state.cachedCourses);

    if (!courseId) {
      return { type: "error", output: `Course '${query}' not found.` };
    }

    try {
      await deleteCourse(courseId);
      state.cachedCourses = await getCourses(); // Refresh cache
      return { type: "success", output: `Course '${query}' deleted.` };
    } catch (err: any) {
      return { type: "error", output: `Failed: ${err.message}` };
    }
  },
};

const mklesson: CommandDefinition = {
  description: "Create a lesson (Admin)",
  execute: async (args) => {
    // Usage: mklesson "Go Mastery" 1 "Intro" "Hello World"
    if (args.length < 4)
      return {
        type: "error",
        output:
          'Usage: mklesson <course_name> <position> "<Title>" "<Content>"',
      };

    const [courseQuery, posStr, title, content] = args;
    const position = parseInt(posStr);

    if (isNaN(position))
      return { type: "error", output: "Position must be a number." };

    // Resolve Course Name -> ID
    // We need to ensure cache exists, just like rmcourse
    if (state.cachedCourses.length === 0) {
      try {
        state.cachedCourses = await getCourses();
      } catch (e) {}
    }
    const courseId = resolveId(courseQuery, state.cachedCourses);

    if (!courseId)
      return { type: "error", output: `Course '${courseQuery}' not found.` };

    try {
      await createLesson(courseId, title, content, position);
      return {
        type: "success",
        output: `Lesson "${title}" created in '${courseQuery}'.`,
      };
    } catch (err: any) {
      return { type: "error", output: `Failed: ${err.message}` };
    }
  },
};

const rmlesson: CommandDefinition = {
  description: "Delete a lesson by Name or ID (Admin)",
  execute: async (args) => {
    if (args.length < 1)
      return { type: "error", output: "Usage: rmlesson <lesson_name_or_id>" };

    const query = args.join(" ");

    // Resolve Lesson Name -> ID
    const lessonId = resolveId(query, state.cachedLessons);

    if (!lessonId) {
      return {
        type: "error",
        output: `Lesson '${query}' not found in current list.\n(Did you run 'lessons <course>' first?)`,
      };
    }

    try {
      await deleteLesson(lessonId);
      state.cachedLessons = state.cachedLessons.filter(
        (l) => l.id !== lessonId,
      );
      return { type: "success", output: `Lesson '${query}' deleted.` };
    } catch (err: any) {
      return { type: "error", output: `Failed: ${err.message}` };
    }
  },
};

// --- EXPORT COMMANDS MAP ---
export const commands: Record<string, CommandDefinition> = {
  help,
  clear,
  register,
  login,
  logout,
  whoami,
  token,
  courses,
  lessons,
  start,
  mkcourse,
  rmcourse,
  mklesson,
  rmlesson,
};
