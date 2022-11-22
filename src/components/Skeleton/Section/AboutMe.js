import React from 'react'

function AboutMe() {

  const languages = [
    ["HTML-2012-E34F26", "html5"],
    ["CSS-2012-1572B6", "css3"],
    ["JavaScript-2015-c5b218", "javascript"],
    ["Python-2018-3776AB", "python"],
    ["SQLite-2018-003B57", "sqlite"],
    ["JSON-2019-333333", "json"],
    ["NumPy-2021-013243","numpy"],
    ["Pandas-2021-150458", "pandas"],
    ["MatPlotLib-2021-E10098","graphql"],
    ["Qt-2020-3FCE51", "qt"],
    ["Markdown-2020-000000", "markdown"],
    ["Bash-2021-4EAA25", "gnubash"],
    ["Django-2021-092E20", "django"],
    ["Awk-2022-666666", "textpattern"],
    ["React-2022-61DAFB", "react"],
    ["Bootstrap-2022-7952B3", "bootstrap"],
    ["Tailwind_CSS-2023-06B6D4", "tailwindcss"]
  ]

  const programs = [
    ["GNU_Linux-0D597F", "linux"],
    ["Git-F05032", "git"],
    ["GIMP-5C5543", "gimp"],
    ["Vim-019733", "vim"],
    ["Neovim-57A143", "neovim"],
    ["Nano-4A90E2", "nano"],
    ["Kdenlive-527EB2", "kdenlive"],
    ["Audacity-0000CC", "audacity"],
    ["MuseScore-1A70B8", "musescore"],
    ["Only_Office-44444444", "onlyoffice"],
    ["Libre_Office-18A303", "libreoffice"],
    ["Open_Office-0E85CD", "apacheopenoffice"],
    ["Microsoft_Office-D83B01", "microsoftoffice"]
  ]

  return <>
    <h1>About Me</h1>
    <p>I study at education department and I want to integrate my coding skills
      to education.  My aim is to create digital and modern education materials.
      The reason why is education in schools still looks same as 100 years ago.
      We just added smartboards. But this cannot be helpful to make students
      smart. So I believe we, who are programmers, should lend a helping hand to
      this situation.</p>
    <p>I am also an open source advocate. I use Linux as my main operating
      system. I publish educational videos about Linux and open source projects on
      my Youtube channel named <a
        href="https://www.youtube.com/@herkesicinlinux">Linux For Everyone</a>.</p>
    <h2>Digital Skill</h2>
    <h3>Programming, Scripting, Declerative, Markup Languages & Modules, Frameworks</h3>
    <div className="flex wrap just-center">
      {languages.map((language, i) => (
        <img key={i} className="margin-small rounded" src={"https://img.shields.io/badge/" + language[0] + "?logo=" + language[1] + "&logoColor=white&style=for-the-badge"} />
      ))}
    </div>
    <h3>Package Programs & OSes</h3>
    <div className="flex wrap just-center">
      {programs.map((program, i) => (
        <img key={i} className="margin-small rounded" src={"https://img.shields.io/badge/" + program[0] + "?logo=" + program[1] + "&logoColor=white&style=for-the-badge"} />
      ))}
    </div>
    <h1>Language Skills</h1>
    <div className="flex align-center">
      Mother Tongue: <img className="margin-small rounded" alt="Turkish" src="https://img.shields.io/badge/Turkish-db0a16?logo=homeassistantcommunitystore&logoColor=white&style=for-the-badge" />
    </div>
    <div className="flex align-center">
      Other Languages: <img className="margin-small rounded" alt="English" src="https://img.shields.io/badge/English-11145b?logo=googleearth&logoColor=white&style=for-the-badge" />
    </div> 
  </>
}

export default AboutMe
