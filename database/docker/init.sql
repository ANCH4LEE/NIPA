--database
CREATE TYPE ticket_status AS ENUM ('Pending', 'Accepted', 'Resolved', 'Rejected');

CREATE TABLE IF NOT EXISTS tickets (
    ticket_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    contact TEXT NOT NULL, --ข้อมูลของคนแจ้ง
    status ticket_status DEFAULT 'Pending',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


--update_at จะอัปเดต ใช้พวก trigger
CREATE OR REPLACE FUNCTION update_ticket()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER update_ticket_trigger
BEFORE UPDATE ON tickets
FOR EACH ROW
EXECUTE FUNCTION update_ticket();

--ไม่ให้ลบ ticket
CREATE RULE prevent_ticket_deletion AS
ON DELETE TO tickets
DO INSTEAD NOTHING;

--ค้นหาสถานะ sort
CREATE INDEX idx_tickets_status ON tickets(status);
CREATE INDEX idx_tickets_updated ON tickets(updated_at);

INSERT INTO tickets (title, description, contact, status)
VALUES
    ('ระบบล่ม', 'เซิร์ฟเวอร์ล่มทั้งระบบ ไม่สามารถใช้งานได้', 'user1@example.com', 'Pending'),
    ('เครื่องพิมพ์เสีย', 'เครื่องพิมพ์ไม่สามารถพิมพ์ได้', 'support@example.com', 'Accepted'),
    ('อินเทอร์เน็ตขัดข้อง', 'การเชื่อมต่ออินเทอร์เน็ตไม่เสถียร', '0123456789', 'Resolved'),
    ('ไฟฟ้าดับ', 'ไฟฟ้าดับในพื้นที่ทำงาน', 'hello@example.com', 'Rejected');




