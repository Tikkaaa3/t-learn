import type { CommandDefinition } from "../types";
import type { Course, Lesson } from "../api/content";
import { loginUser, generateApiKey, registerUser, getUserStats } from "../api/auth";
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
  path: string[]; // ["courses"] or ["courses", "Rust Basics"]
  currentCourse: Course | null; // Current course we're in
  cachedCourses: Course[];
  cachedLessons: Lesson[];
}

// Initial State
const state: ShellState = {
  user: localStorage.getItem("t_learn_user") || null,
  path: [],
  currentCourse: null,
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
  if (state.path.length === 0) {
    return "$";
  }
  return `${state.path.join("/")} $`;
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
  whoami                        - Show current user and stats
  token                         - Generate CLI API Key
  ls [courses]                  - List courses or lessons
  cd <course_name>              - Enter a course
  cd ..                         - Go back to courses root
  start <lesson_name>           - Start a lesson task
\`\`\`
`,
    };
  },
};

const clear: CommandDefinition = {
  description: "Clear terminal",
  execute: async () => {
    return { type: "info", output: "" };
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
    state.path = [];
    state.currentCourse = null;
    state.cachedCourses = [];
    state.cachedLessons = [];

    return { type: "success", output: "Logged out successfully." };
  },
};

const whoami: CommandDefinition = {
  description: "Show current user and stats",
  execute: async () => {
    if (!state.user) {
      return { type: "error", output: "Not logged in." };
    }

    try {
      const stats = await getUserStats();
      return {
        type: "success",
        output: `**${stats.username}**\nCompleted tasks: ${stats.completed_tasks}`,
      };
    } catch (err: any) {
      // Fallback if stats endpoint fails
      return { type: "success", output: state.user };
    }
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

const ls: CommandDefinition = {
  description: "List courses or lessons",
  execute: async () => {
    // If at root, show hint
    if (state.path.length === 0) {
      return {
        type: "info",
        output: "Type 'cd courses' to browse courses.",
      };
    }

    // If we're in courses directory (not in a specific course)
    if (state.path.length === 1 && state.path[0] === "courses") {
      try {
        const courses = await getCourses();
        state.cachedCourses = courses;

        if (courses.length === 0)
          return { type: "info", output: "No courses found." };

        const list = courses
          .map((c) => `📁 **${c.title}**/`)
          .join("\n");

        return { type: "info", output: `### Courses:\n${list}` };
      } catch (err: any) {
        return {
          type: "error",
          output: `Failed to fetch courses: ${err.message}`,
        };
      }
    }

    // If we're inside a course, list lessons
    if (state.currentCourse) {
      try {
        const lessons = await getLessons(state.currentCourse.id);
        state.cachedLessons = lessons;

        if (lessons.length === 0)
          return { type: "info", output: `No lessons in this course.` };

        const list = lessons
          .map((l) => {
            const mark = l.completed ? "✓" : " ";
            return `[${mark}] 📄 ${l.title}`;
          })
          .join("\n");

        return {
          type: "info",
          output: `### Lessons in ${state.currentCourse.title}:\n${list}`,
        };
      } catch (err: any) {
        return { type: "error", output: `Failed: ${err.message}` };
      }
    }

    return { type: "info", output: "Nothing to list here." };
  },
};

const cd: CommandDefinition = {
  description: "Change directory to a course",
  execute: async (args) => {
    if (args.length === 0) {
      return { type: "error", output: "Usage: cd <directory> or cd .." };
    }

    // cd .. goes back one level
    if (args[0] === "..") {
      if (state.path.length === 0) {
        return { type: "info", output: "Already at root." };
      }
      
      // If we're in a course, go back to courses directory
      if (state.path.length === 2) {
        state.path = ["courses"];
        state.currentCourse = null;
        state.cachedLessons = [];
        return { type: "info", output: "Moved to courses directory." };
      }
      
      // If we're in courses directory, go back to root
      if (state.path.length === 1) {
        state.path = [];
        return { type: "info", output: "Moved to root." };
      }
      
      return { type: "info", output: "Moved back." };
    }

    const target = args.join(" ");

    // cd courses (from root)
    if (target === "courses" && state.path.length === 0) {
      state.path = ["courses"];
      return { type: "info", output: "Moved to courses directory. Type 'ls' to see courses." };
    }

    // cd <course_name> (from courses directory)
    if (state.path.length === 1 && state.path[0] === "courses") {
      // Ensure cache
      if (state.cachedCourses.length === 0) {
        try {
          state.cachedCourses = await getCourses();
        } catch (e) {}
      }

      const courseId = resolveId(target, state.cachedCourses);
      if (!courseId)
        return { type: "error", output: `Course '${target}' not found.` };

      const course = state.cachedCourses.find((c) => c.id === courseId);
      if (!course)
        return { type: "error", output: `Course '${target}' not found.` };

      state.currentCourse = course;
      state.path = ["courses", course.title];

      // Fetch lessons immediately
      try {
        const lessons = await getLessons(courseId);
        state.cachedLessons = lessons;
      } catch (e) {}

      return {
        type: "info",
        output: `Entered course: **${course.title}**\nType 'ls' to see lessons.`,
      };
    }

    return { type: "error", output: `Cannot cd to '${target}' from current location.` };
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
  ls,
  cd,
  start,
  mkcourse,
  rmcourse,
  mklesson,
  rmlesson,
};
