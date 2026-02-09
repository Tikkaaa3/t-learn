import type { CommandDefinition } from "../types";

// Placeholder Functions (I will connect these to API later)

const help: CommandDefinition = {
  description: "List available commands",
  execute: async () => {
    return {
      type: "info",
      output: `Available commands:
  help      - Show this message
  clear     - Clear the terminal screen
  login     - Log in to the platform
  courses   - List available courses
  whoami    - Show current user`,
    };
  },
};

const login: CommandDefinition = {
  description: "Log in to the platform",
  execute: async (args) => {
    if (args.length < 2) {
      return { type: "error", output: "Usage: login <username> <password>" };
    }
    // TODO: Connect to Real Backend
    return {
      type: "success",
      output: `Mock login successful for user: ${args[0]}`,
    };
  },
};

const courses: CommandDefinition = {
  description: "List available courses",
  execute: async () => {
    // TODO: Fetch from Backend
    return {
      type: "info",
      output: `
ID                                    | TITLE
--------------------------------------+----------------------
550e8400-e29b-41d4-a716-446655440000 | Go HTTP Mastery
770e8400-e29b-41d4-a716-446655441111 | Advanced SQL`,
    };
  },
};

const whoami: CommandDefinition = {
  description: "Show current user",
  execute: async () => {
    return { type: "info", output: "guest" };
  },
};

// The Registry Map

export const commands: Record<string, CommandDefinition> = {
  help,
  login,
  courses,
  whoami,
  // 'clear' is handled specially in the hook
};
