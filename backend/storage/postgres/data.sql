delete from smartphones;
SELECT setval(pg_get_serial_sequence('smartphones', 'id'), coalesce(max(id),0) + 1, false) FROM smartphones;
insert into smartphones (model, producer, memory, ram, display_size, description, image_path, price)
values (
    'iPhone 16',
    'Apple',
    128,
    8,
    6.1,
   'Introducing the all-new iPhone 16 where innovation meets elegance. With a sleek design and cutting-edge technology, the iPhone 16 delivers a stunning display, incredible camera capabilities, and lightning-fast performance that transforms the way you experience mobile devices. Whether you are capturing memories in breathtaking detail, enjoying your favorite content in vibrant color, or seamlessly multitasking, the iPhone 16 elevates every moment. Powered by Apple most advanced chipset, it brings unmatched speed, efficiency, and security to your fingertips. Step into the future with the iPhone 16 - a new standard in smartphone excellence.',
    'https://c.dns-shop.ru/thumb/st1/fit/0/0/1043f341d851923dda2ac92e50f089a1/14ce8c6a5fbaef30feb3cb6b7d742546c045c44eb9207be4acec68cade72a7cf.jpg.webp',
    999
),
(
    'C61',
    'POCO',
    64,
    3,
    6.78,
    'The POCO C61 64GB smartphone comes in a black plastic body with a glass back panel. Corning Gorilla Glass coating protects the screen from scuffs and scratches. The 6.78 IPS display shows a clear and bright image with a resolution of 1600x720 dpi. The 90Hz refresh rate is sufficient for smooth animated effects.
The POCO C61 with 64GB of storage, an 8-core MediaTek processor, and 3GB of RAM is suitable for calls and social media without installing games or large applications. There a separate slot for a microSD external drive up to 1TB, distinct from the SIM card slots. The 5000 mAh battery holds a charge throughout the day. The smartphone is equipped with universal USB Type-C connectors for charging and a 3.5mm jack for headphones. It has a fingerprint scanner built into the power button and a face recognition feature in the camera.',
    'https://c.dns-shop.ru/thumb/st1/fit/0/0/a3df61f4a11c4524eda958781c36b1b3/4597932c0d86bff8db78fc5242cdeefade448b1ede21de349fdf2202bebdeaf8.jpg.webp',
    59

),
(
    'Galaxy A55',
    'Samsung',
    128,
    8,
    6.6,
    'The Samsung Galaxy A55 6.6” 5G 128GB smartphone in blue features a 6.6-inch Super AMOLED display. An adaptive 120Hz refresh rate ensures tear-free image output and smooth video playback. The Super AMOLED panel provides deep blacks, making the image sharp with natural color reproduction. The model is equipped with 3 rear cameras. The 50MP main camera supports automatic stabilization technology, taking clear shots in motion. The macro camera can capture small objects, while the ultra-wide-angle lens allows many objects to fit into the frame at once.
The Samsung Galaxy A55 smartphone is powered by an 8-core Exynos 1480 processor, in which graphics performance has been increased by 32% compared to the previous generation — the device works without delays in multitasking mode and can run resource-intensive games. The cooling system features a large-area vapor chamber, quickly dissipating heat and preventing overheating.
The body is made of metal, protected to IP67 standard, preventing dust and water from entering even with direct contact. The screen is protected by Corning Gorilla Glass Victus+, which is scratch-resistant.',
    'https://c.dns-shop.ru/thumb/st1/fit/0/0/9415d7ccfb69d3454721c2a8ecc7167c/610c3f8c51503fbd4f2683bb886922c50fc4a9be8a808e14857bc28a3a061046.jpg.webp',
    299
),
(
    'Redmi Note 14 Pro',
    'Xiaomi',
    512,
    12,
    6.67,
    'The Xiaomi Redmi Note 14 Pro is a powerful and stylish smartphone with a large screen and high resolution. The device features a bright and clear AMOLED display with a 6.67-inch diagonal and 2400x1080 pixel resolution, which ensures excellent image quality and comfortable use.
The smartphone runs on the Android 14 operating system with the proprietary HyperOS shell, offering a convenient and intuitive interface.
The Xiaomi Redmi Note 14 Pro is equipped with a powerful MediaTek Helio G100 Ultra processor, ensuring fast and smooth operation of the device. The smartphone also has a large amount of LPDDR4X RAM, allowing it to easily handle multitasking and run demanding applications.
The triple main camera with optical stabilization allows for high-quality and vibrant photos and videos. The device also supports artificial intelligence features that help improve the quality of images and videos.
The smartphone has a capacious 5500 mAh battery, providing long operating time without recharging. The Xiaomi Redmi Note 14 Pro also supports 45W fast charging, allowing for quick battery replenishment.',
    'https://c.dns-shop.ru/thumb/st1/fit/500/500/468d30e0c1de5961e891c391dca5f587/a680bc390ab900b2efcd1b086ce40b2f98f819e0eafa2722817885201d6fb7d7.jpg.webp',
    379
),
(
    'Moto G55',
    'Motorola',
    128,
    4,
    6.49,
    'The Motorola Moto G55 is a powerful and stylish smartphone with a large screen and high resolution. The device features a bright and clear IPS display with a 6.49-inch diagonal and 2400x1080 pixel resolution, which ensures excellent image quality and comfortable use.
The smartphone runs on the Android 14 operating system with the proprietary My UX shell, offering a convenient and intuitive interface.
The Motorola Moto G55 is equipped with a powerful MediaTek Dimensity 7025 processor, ensuring fast and smooth operation of the device. The smartphone also has a large amount of LPDDR4X RAM, allowing it to easily handle multitasking and run demanding applications.
The dual main camera with optical stabilization allows for high-quality and vibrant photos and videos.
The smartphone has a capacious 5000 mAh battery, providing long operating time without recharging. The Motorola Moto G55 also supports 30W fast charging, allowing for quick battery replenishment.',
    'https://c.dns-shop.ru/thumb/st1/fit/500/500/f80c91b992ec99418ed5b490de58996d/818869be182d74294ac028a3e8fc441a2d528b76a9a93aea673b68a8b66d0608.jpg.webp',
    179
),
(
    'Pixel 8a',
    'Google',
    128,
    8,
    6.1,
    'Magic photo and video editing. Google AI will help you achieve your goals. All-day battery and durable design. 7 years of feature and security updates. The same chip as in the Pixel 8 Pro.',
    'https://c.dns-shop.ru/thumb/st1/fit/500/500/fcdb06077c903a6c1aae076d92e9ad15/44632b916e2df553f76c8c0481cd03c0b030ece24f352d7034fc89fa772f463c.jpg.webp',
    499
),
(
    'ROG Phone 9 Pro',
    'ASUS',
    512,
    16,
    6.78,
    'The ASUS ROG Phone 9 Pro represents an outstanding achievement in mobile devices, combining advanced gaming capabilities and a gaming design. The smartphone boasts an IP68 protection rating, ensuring reliable defense against dust and water. This allows users to game anytime, anywhere, without worrying about potential damage.
The ASUS ROG Phone 9 Pro is powered by the 3-nanometer Snapdragon 8 Elite mobile platform, which provides high performance and seamless PC synchronization.
The smartphone features an advanced 6.78-inch AMOLED display with an adaptive refresh rate that minimizes power consumption. The ROG Phone 9 supports a refresh rate of up to 185Hz via Game Genie, ensuring an ultra-responsive 720Hz touch sampling rate. The display peak brightness is 2500 nits, guaranteeing a clear image even in bright sunlight. Delta-E < 1 color accuracy ensures high-quality visuals.
The ASUS ROG Phone 9 Pro utilizes advanced technology that increases battery energy density by 3.5%. This allows for an increased battery capacity of up to 5800 mAh, significantly improving the device overall operating time. Furthermore, the battery is capable of withstanding over 1000 charge and discharge cycles while retaining more than 80% of its original capacity.',
    'https://c.dns-shop.ru/thumb/st1/fit/0/0/b250dd1f13be1e5492f1c00112efcca5/630e0adaf67ae3adb97c24ac889fff56c5e45963195129c10e54faa8fcf94606.jpg.webp',
    1199
),
(
    'Phone 2a',
    'Nothing',
    128,
    8,
    6.7,
    'The Nothing Phone 2a features an unusual design, a semi-transparent body with a horizontal arrangement of two rear camera modules (a 50MP main sensor and a 50MP ultrawide) and three LEDs. The smartphone edges are flat, as is the 6.7-inch AMOLED screen. The display resolution is FHD+, with a 120Hz refresh rate and a peak brightness of 1300 nits. The front camera is 32MP.

The hardware is based on the MediaTek Dimensity 7200 Pro chip. The battery capacity is 5000 mAh, with 45W wired fast charging available. The Nothing Phone (2a) is protected according to the IP54 standard. The OS is Android 14 with the proprietary Nothing OS 2.5 interface.',
    'https://c.dns-shop.ru/thumb/st4/fit/0/0/85d58e56e0c912b5203bf63cd235970b/95b0b2dfb651b1a89ca7cb57f86e9640138ab149ae9f986d3782efdc65b8905f.jpg.webp',
    299
);

delete from carts;
SELECT setval(pg_get_serial_sequence('carts', 'id'), coalesce(max(id),0) + 1, false) FROM carts;

delete from cart_items;
SELECT setval(pg_get_serial_sequence('cart_items', 'id'), coalesce(max(id),0) + 1, false) FROM cart_items;

delete from users;
SELECT setval(pg_get_serial_sequence('users', 'id'), coalesce(max(id),0) + 1, false) FROM users;
insert into users (name, password, role)
values
    ('admin', '$2a$10$4ntk1IwmZXCKR/QGp9cZU.2JVvqYWwM9uyxKlC7pPR5suRlf4Bkx.', 'admin'),
    ('user1', '$2a$10$f9sCd/oZg9GriPoHVrHMT.4KIr6dOwmQbU5FDCQdYgxYm3Xc6pQqa', 'user'),
    ('user2', '$2a$10$FKAxhcBTV9/yZbNk9OhbpeDrW5RMSrNFKT8w1OHGENR.sV.kqgUEi', 'user');

insert into cart_items (cart_id, smartphone_id)
values
    (1, 1),
    (1, 3),
    (2, 2);

delete from reviews;
SELECT setval(pg_get_serial_sequence('reviews', 'id'), coalesce(max(id),0) + 1, false) FROM reviews;
insert into reviews (smartphone_id, user_id, rating, comment)
values
    (1, 2, 5, 'Best as always.'),
    (1, 3, 4, null),
    (3, 2, 1, 'Stopped working just after 2 days :('),
    (4, 3, 3, null),
    (6, 2, 5, 'Great camera, powerful CPU'),
    (7, 3, 5, '16 gb of RAM is absolutely insane');
