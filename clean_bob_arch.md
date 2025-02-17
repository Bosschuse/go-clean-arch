# Summary
Clean bob arch เป็น framework ตัวหนึ่งที่ใช้ในการทำงาน เพื่อให้ code สะอาด และ ในทีม ทำไปในแนวทางเดียวกัน โดย โมดูลแต่ละส่วนสามารถทำงานแยกกันโดยอิสระ ไม่ขึ้นต่อกันทำให้ง่ายต่อการแก้ไขโดยไม่สร้างผลกระทบต่อส่วนอื่น ๆ โดยบังเอิญ
## Profit
Test able code
ทุกโมดูลแยกออกจากกันอย่างอิสระ 

ทำให้สามารถแก้ไขโดยไม่กระทบส่วนอื่น
## Core Idea
### Entities
กำหนดโครงสร้างของข้อมูล เช่น Interface, struct
### Usecase
bussiness logic ใด ๆ ก่อนจะเข้าถึง database
### Controller
รับส่งข้อมูล http จาก client
### Repositories
ส่วนของการเชื่อมต่อ database เพื่อ query ข้อมูล

