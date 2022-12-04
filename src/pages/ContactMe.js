import Section from "../components/Skeleton/Section";
import { LanguageContext } from "../contexts/LanguageContext"
import { useContext } from "react";
import { directs, h2s, mainTitle, others } from "../translations/Contact";

function ContactMe() {

  const { language } = useContext(LanguageContext)

  return <Section>
    {mainTitle[language]}
    {h2s[language][0]}
    <div className="flex flex-col gap-3">
      {directs[language]}
    </div>
    {h2s[language][1]}
    <div className="flex flex-col gap-3">
      {others[language]}
    </div>
  </Section>
}

export default ContactMe