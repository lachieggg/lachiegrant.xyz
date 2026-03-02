import os
from pypdf import PdfWriter, PdfReader

reader = PdfReader("2026-cv.pdf")
writer = PdfWriter()

for page in reader.pages:
    writer.add_page(page)

writer.add_metadata({"/Title": "Resume"})

with open("2026-cv-out.pdf", "wb") as f:
    writer.write(f)

print("Done! Saved as 2026-cv-out.pdf")


