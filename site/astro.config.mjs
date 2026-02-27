import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";

export default defineConfig({
  site: "https://cloverhoundps.github.io",
  base: process.env.BASE_URL || "/epic-fhir-cli",
  integrations: [
    starlight({
      title: "fhir-cli",
      social: [
        {
          icon: "github",
          label: "GitHub",
          href: "https://github.com/CloverhoundPS/epic-fhir-cli",
        },
      ],
      sidebar: [
        {
          label: "Getting Started",
          items: [
            { label: "Installation", slug: "installation" },
            { label: "Quick Start", slug: "getting-started" },
            { label: "Configuration", slug: "configuration" },
          ],
        },
        {
          label: "Authentication",
          items: [
            { label: "Auth Setup", slug: "authentication" },
            { label: "Troubleshooting", slug: "auth-troubleshooting" },
          ],
        },
        {
          label: "Clinical Resources",
          items: [
            { label: "Patient", slug: "patient" },
            { label: "Observation", slug: "observation" },
            { label: "Condition", slug: "condition" },
            { label: "Medication", slug: "medication" },
            { label: "Allergy", slug: "allergy" },
            { label: "Encounter", slug: "encounter" },
            { label: "Procedure", slug: "procedure" },
            { label: "Diagnostic Report", slug: "diagnostic" },
            { label: "Immunization", slug: "immunization" },
            { label: "Care Plan", slug: "careplan" },
            { label: "Document", slug: "document" },
            { label: "Family History", slug: "family" },
          ],
        },
        {
          label: "Administrative Resources",
          items: [
            { label: "Appointment", slug: "appointment" },
            { label: "Coverage", slug: "coverage" },
            { label: "Practitioner", slug: "practitioner" },
          ],
        },
        {
          label: "Advanced",
          items: [
            { label: "Generic Resource Commands", slug: "resource" },
            { label: "Output Formats", slug: "output-formats" },
            { label: "Epic Sandbox", slug: "epic-sandbox" },
          ],
        },
      ],
    }),
  ],
});
