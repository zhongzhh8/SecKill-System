import csv

# 200普通用户，5商家，每个商家10种券，每种券20张，应该全部抢完

customer_num = 200
customer_register_path = "customer_register.csv"
customer_seckill_path = "customer_seckill.csv"

saler_num = 5
coupon_num = 10
saler_register_path = "saler_register.csv"
saler_create_coupon_path = "saler_create_coupon.csv"


# 单纯用户注册
def create_customer_register_csv():
    with open(customer_register_path, 'w') as f:
        csv_write = csv.writer(f)
        for id in range(customer_num):
                    username = 'customer_' + str(id)
                    password = '123456'
                    csv_write.writerow([username, password])


# 用户登录，秒杀所有商家优惠券
def create_customer_seckill_csv():
    with open(customer_seckill_path, 'w') as f:
        csv_write = csv.writer(f)
        for id in range(customer_num):
            for saler_id in range(saler_num):
                for coupon_id in range(coupon_num):
                    username = 'customer_' + str(id)
                    password = '123456'
                    saler = 'saler_' + str(saler_id)
                    coupon_id = 'saler_' + str(saler_id) + '_'+ str(coupon_id)
                    csv_write.writerow([username, password, saler, coupon_id])

# 单纯商家注册
def create_saler_register_csv():
    with open(saler_register_path, 'w') as f:
        csv_write = csv.writer(f)
        for saler_id in range(saler_num):
                    username = 'saler_' + str(saler_id)
                    password = '123456'
                    csv_write.writerow([username, password])

# 商家登录，创建优惠券
def create_saler_create_coupon_csv():
    with open(saler_create_coupon_path, 'w') as f:
        csv_write = csv.writer(f)
        for saler_id in range(saler_num):
            for coupon_id in range(coupon_num):
                username = 'saler_' + str(saler_id)
                password = '123456'
                coupon_id = 'saler_' + str(saler_id) + '_'+ str(coupon_id)
                csv_write.writerow([username, password, coupon_id])


if __name__ == "__main__":
    create_customer_register_csv()
    create_customer_seckill_csv()
    create_saler_register_csv()
    create_saler_create_coupon_csv()