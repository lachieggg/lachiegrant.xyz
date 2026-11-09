import os
from pypdf import PdfWriter, PdfReader

ext = ".pdf"
fname = "2026-cv"

reader = PdfReader(fname + ext)
writer = PdfWriter()

for page in reader.pages:
    writer.add_page(page)

writer.add_metadata({"/Title": "Resume"})

with open(fname + "-out" + ext, "wb") as f:
    writer.write(f)

print("Done! Saved as " + fname + "-out" + ext)