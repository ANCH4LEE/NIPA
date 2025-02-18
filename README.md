# 🔌วิธีติดตั้งและใช้งาน
- ทำการเปิด Folder database จาก [Database](database)
- จะต้องเพิ่มไฟล์ .env และใส่ข้อมูลจากด้านล่าง
```bash 
# PostgreSQL Environment Variables
POSTGRES_DB=ticket
POSTGRES_USER=ticket_user
POSTGRES_PASSWORD=ticket_password
POSTGRES_PORT=5432

# pgAdmin Environment Variables
PGADMIN_DEFAULT_EMAIL=admin123@gmail.com
PGADMIN_DEFAULT_PASSWORD=password
PGADMIN_PORT=5055
```

- กด save File จากนั้น Run Folder ด้วย Docker ใช้คำสั่ง
```bash
docker-compose build
docker-compose up -d 
```

- ทำการเปิด Folder back จาก [Back](Back)
- จะต้องเพิ่มไฟล์ .env และใส่ข้อมูลจากด้านล่าง
```bash 
#เปิด Terminal ใน Vs code หรือ Power shell หรือ Command line
#พิมพ์คำสั่ง ipconfig เพื่อค้นหา IP address ของ IPv4

#จากนั้นคัดลอก IPv4 ไปเปลี่ยนใน POSTGRES_HOST ในไฟล์ .env
# App
APP_PORT=8080

# Database
POSTGRES_HOST=192.168.1.109 
POSTGRES_PORT=5432
POSTGRES_USER=ticket_user
POSTGRES_PASSWORD=ticket_password
POSTGRES_DBNAME=ticket
POSTGRES_SSLMODE=disable
```
- กด save File จากนั้น Run Folder ด้วย Docker ใช้คำสั่ง
```bash
docker-compose build
docker-compose up -d 
```
- สำหรับ PostgreSQL เปิด localhost:5050
- ทำการเปิดโฟลเดอร์ helpdesk จาก [Front](front/helpdesk)
- เพิ่ม Proxy สำหรับเชื่อมต่อกับฝั่ง Backend ในไฟล์ package.json
```bash 
"proxy": "http://localhost:8080"
}
```
- Run Server
```bash 
npm start
```
- สำหรับ React เปิด localhost:3000 
