"""
Custom Claude Code agent with Skills support.

This agent extends the standard ClaudeCode agent to support loading
custom skills from a directory before running the agent.
"""

from pathlib import Path

from harbor.agents.installed.claude_code import ClaudeCode
from harbor.environments.base import BaseEnvironment


class ClaudeCodeWithSkills(ClaudeCode):
    """
    Claude Code agent with custom Skills support.
    
    Skills are loaded from a directory and copied into the container's
    .claude/skills/ directory before running the agent.
    
    Args:
        skill_dir: Path to directory containing skill subdirectories
        skills: Comma-separated list of skill names to include (default: all)
    
    Usage:
        # Load all skills
        --ak "skill_dir=skills"
        
        # Load specific skills
        --ak "skill_dir=skills" --ak "skills=atlas-schema-mgmt"
        
        # Load no skills (baseline)
        --ak "skill_dir=skills" --ak "skills="
    """
    
    def __init__(
        self,
        logs_dir: Path,
        skill_dir: Path | str | None = None,
        skills: str | None = None,
        **kwargs,
    ):
        super().__init__(logs_dir, **kwargs)
        self._skill_dir = Path(skill_dir) if skill_dir else None
        # Parse skills filter: None means all, empty string means none
        if skills is None:
            self._skills_filter = None  # Load all
        elif skills == "":
            self._skills_filter = set()  # Load none
        else:
            self._skills_filter = set(s.strip() for s in skills.split(","))
    
    @staticmethod
    def name() -> str:
        return "claude-code-with-skills"
    
    def _should_load_skill(self, skill_name: str) -> bool:
        """Check if a skill should be loaded based on the filter."""
        if self._skills_filter is None:
            return True  # No filter, load all
        return skill_name in self._skills_filter
    
    async def setup(self, environment: BaseEnvironment) -> None:
        """Setup the agent, including copying skills into the container."""
        # First, run the parent setup (installs Claude Code)
        await super().setup(environment)
        
        # If we have a skill directory, copy skills into the container
        if self._skill_dir and self._skill_dir.exists():
            # Create the skills directory in the workspace
            await environment.exec(command="mkdir -p /workspace/.claude/skills")
            
            loaded_skills = []
            
            # Copy each skill subdirectory
            for skill_path in self._skill_dir.iterdir():
                if skill_path.is_dir() and (skill_path / "SKILL.md").exists():
                    skill_name = skill_path.name
                    
                    # Check if this skill should be loaded
                    if not self._should_load_skill(skill_name):
                        continue
                    
                    target_dir = f"/workspace/.claude/skills/{skill_name}"
                    
                    # Create the skill directory
                    await environment.exec(command=f"mkdir -p {target_dir}")
                    
                    # Upload each file in the skill directory
                    for file_path in skill_path.rglob("*"):
                        if file_path.is_file():
                            relative = file_path.relative_to(skill_path)
                            target_path = f"{target_dir}/{relative}"
                            
                            # Create parent directories if needed
                            parent_dir = str(Path(target_path).parent)
                            if parent_dir != target_dir:
                                await environment.exec(command=f"mkdir -p {parent_dir}")
                            
                            await environment.upload_file(
                                source_path=file_path,
                                target_path=target_path,
                            )
                    
                    loaded_skills.append(skill_name)
            
            if loaded_skills:
                print(f"Loaded skills: {', '.join(loaded_skills)}")
            else:
                print("No skills loaded (baseline)")
            
            # Verify skills were loaded
            result = await environment.exec(command="ls -la /workspace/.claude/skills/")
            print(f"Skills directory: {result.stdout}")
