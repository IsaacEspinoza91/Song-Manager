-- Bloque anonimo
DO $$ 
DECLARE prisioneros_id BIGINT;
DECLARE laferte_id BIGINT;
DECLARE ases_id BIGINT;
DECLARE faceb_id BIGINT;
DECLARE bunny_id BIGINT;
DECLARE easy_id BIGINT;
DECLARE coldplay_id BIGINT;
DECLARE bunkers_id BIGINT;
BEGIN 

    INSERT INTO artists (name, genre, country, bio, image_url) 
    VALUES ('Los Prisioneros', 'Rock Pop', 'Chile', 'Banda de rock chilena de San Miguel, fundamental en la música latinoamericana.', 'https://i.scdn.co/image/ab6761610000e5eb3e37d9994c60ed38760ecb72')
    RETURNING id INTO prisioneros_id;

    -- La voz de los 80 1984
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('La voz de los ''80', '1984-12-13', 'LP', 'https://m.media-amazon.com/images/I/71cVtLcJWyL.jpg') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('La voz de los ''80', 248), ('Brigada de negro', 226), ('Latinoamérica es un pueblo al sur de EE.UU.', 242),
        ('Eve-Evelyn', 264), ('Sexo', 288), ('¿Quién mató a Marilyn?', 188), ('Paramar', 225),
        ('No necesitamos banderas', 309), ('Mentalidad televisiva', 256), ('Nunca quedas mal con nadie', 251)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, prisioneros_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, prisioneros_id FROM current_album;
    

    -- Pateando piedras 1986
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('Pateando Piedras', '1986-09-15', 'LP', 'https://m.media-amazon.com/images/I/71wmK3X0ijL.jpg')
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Muevan las industrias', 248), ('¿Por qué no se van?', 181), ('El baile de los que sobran', 322), 
        ('Estar solo', 273), ('Exijo ser un héroe', 344), ('Quieren dinero', 310), ('Por favor', 212), 
        ('¿Por qué los ricos?', 297), ('Una mujer que no llame la atención', 202), ('Independencia cultural', 274)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, prisioneros_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, prisioneros_id FROM current_album;


    -- La cultura de la basura 1987
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('La cultura de la basura', '1987-12-03', 'LP', 'https://m.media-amazon.com/images/I/61XSG2vrXfL.jpg')
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Somos solo rudio', 81), ('De la cultura de la basura', 187), ('Que no destrocen tu vida', 252), ('Usted y su ambición', 227),
        ('Cuando te vayas', 291), ('Jugar a la guerra', 277), ('Algo tan moderno', 285), ('Maldito sudaca', 136), 
        ('Lo estamos pasando muy bien', 348), ('El es mi ídolo', 258), ('El vals', 187), ('Otro día', 276), ('Pa pa pa', 210), ('Poder elegir', 482)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, prisioneros_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, prisioneros_id FROM current_album;

    -- Corazones (1990)
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('Corazones', '1990-05-20', 'LP', 'https://m.media-amazon.com/images/I/71TI05J2YgL.jpg')
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Tren al sur', 336), ('Amiga mía', 243), ('Con suavidad', 302), ('Corazones rojos', 210), 
        ('Cuéntame una historia original', 232), ('Estrechez de corazón', 384), ('Por amarte', 362), 
        ('Noche en la ciudad', 276), ('Es demasiado triste', 290)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, prisioneros_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, prisioneros_id FROM current_album;




    -- MON LAFERTE
    INSERT INTO artists (name, genre, country, bio, image_url) 
    VALUES ('Mon Laferte', 'Rock Pop', 'Chile', 'Norma Monserrat Bustamante Laferte, más conocida por su nombre artístico Mon Laferte, es una cantante y compositora chilenomexicana', 'https://i.scdn.co/image/ab6761610000e5ebd4d9941b20bbf2bbe3714acc')
    RETURNING id INTO laferte_id;

    -- Mon Laferte Vol.1/Edicion especial (2015)
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('Mon Laferte (Vol.1/Edición Especial)', '2015-01-31', 'LP', 'https://i.scdn.co/image/ab67616d000082c17a9c61e7c19a07afe65a3c9d')
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Tormento', 276), ('El cristal', 181), ('El diablo', 248), ('La visita', 227), ('Amor completo', 241),
        ('Un alma en pena', 142), ('Tu falta de querer', 278), ('Salvador', 201), ('Si tú me quisieras', 202),
        ('Malagradecido', 187), ('La noche del día que llovió en verano', 58), ('Bonita', 219), ('Orgasmo para dos - En vivo', 248),
        ('Vuelve por favor - En vivo', 357), ('Flor de amapola - En vivo', 259), ('Igual que yo - En vivo', 347),
        ('Si tú me quisieras - En vivo', 211), ('Tu falta de querer', 278)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, laferte_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, laferte_id FROM current_album;


    -- Femme Fatale (2025)
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('Femme Fatale', '2025-10-24', 'LP', 'https://i.scdn.co/image/ab67616d000082c1138aa95e930672e18daf4f66')
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Femme Fatale', 250), ('Mi hombre', 198), ('Otra noche de llorar', 219), ('Esto es amor', 273), 
        ('Veracruz', 262), ('El gran señor', 203), ('Las flores que dejaste en la mesa', 209), ('1:30', 160),  
        ('La tirana', 244), ('Hasta que nos despierte la soledad', 227), ('Melancolía', 230), 
        ('Ocupa mi piel', 225), ('My One And Only Love', 181), ('Vida normal', 249)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, laferte_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, laferte_id FROM current_album;








    -- BAD BUNNY
    INSERT INTO artists (name, genre, country, bio, image_url) 
    VALUES ('Bad Bunny', 'Pop', 'Puerto Rico', 'Benito Antonio Martínez Ocasio, conocido artísticamente como Bad Bunny, es un cantante, compositor y productor discográfico puertorriqueño.', 'https://i.scdn.co/image/ab6761610000e5eb81f47f44084e0a09b5f0fa13')
    RETURNING id INTO bunny_id;

    -- X100PRE
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('X100PRE', '2018-12-24', 'LP', 'https://m.media-amazon.com/images/I/719-bcxT-xL.jpg') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('NI BIEN NI MAL', 236), ('200 Mph', 170), ('¿Quien Tu Eres?', 159), ('Caro', 229), ('Tenemos Que Hablar', 224),
        ('Otra Noche en Miami', 233), ('Ser Bichote', 193), ('Si Estuviésemos Juntos', 169), ('Solo de Mi', 197), 
        ('Cuando Perriabas', 188), ('La Romana', 300), ('Como Antes', 230), ('RLNDT', 284), ('Estamos Bien', 208), ('MIA', 210)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, bunny_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, bunny_id FROM current_album;

    -- OASIS
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('OASIS', '2019-06-28', 'LP', 'https://i.scdn.co/image/ab67616d000082c14891d9b25d8919448388f3bb') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('MOJAITA', 187), ('YO LE LLEGO', 249), ('CUIDAD POR AHÍ', 198), ('QUE PRETENDES', 222), 
        ('LA CANCIÓN', 242), ('UN PESO', 277), ('ODIO', 270), ('COMO UN BEBÉ', 218)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, bunny_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, bunny_id FROM current_album;

    -- YHLQMDLG
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('YHLQMDLG', '2020-02-29', 'LP', 'https://m.media-amazon.com/images/I/81RqzDbsNBL.jpg') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Si Veo a Tu Mamá', 170), ('La Difícil', 163), ('Pero Ya No', 160), ('La Santa', 206), ('Yo Perreo Sola', 172),
        ('Bichiyal', 196), ('Soliá', 159), ('La Zona', 136), ('Que Malo', 167), ('Vete', 192), ('Ignorantes', 210),
        ('A Tu Merced', 175), ('Una Vez', 232), ('Safaera', 295), ('25/8', 243), ('Está Cabrón Ser Yo', 227), 
        ('Puesto Pa Guerrial', 190), ('P FKN R', 258), ('Hablamos Mañana', 240), ('<3', 157)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, bunny_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, bunny_id FROM current_album;

    -- LAS QUE NO IBAN A SALIR
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('LAS QUE NO IBAN A SALIR', '2020-05-10', 'LP', 'https://m.media-amazon.com/images/I/61ff3vX88nL.jpg') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('SI ELLA SALE', 143), ('MÁS DE UNA CITA', 183), ('BYE ME FUI', 178), ('CANCIÓN CON YANDEL', 209), ('PA ROMPERLA', 194),
        ('BAD CON NICKY', 202), ('BENDICIONES', 155), ('CÓMO SE SIENTE - Remix', 227), ('RONCA FREESTYLE', 150), ('EN CASITA', 176)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, bunny_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, bunny_id FROM current_album;

    -- EL ÚLTIMO TOUR DEL MUNDO
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('EL ÚLTIMO TOUR DEL MUNDO', '2020-11-27', 'LP', 'https://m.media-amazon.com/images/I/81fHDxOVd4L.jpg') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('El mundo es mío', 165), ('Te mudaste', 130), ('Hoy cobré', 162), ('Maldita pobreza', 213), 
        ('La noche de anoche', 203), ('Te deseo lo mejor', 139), ('Yo visto así', 191), ('Haciendo que me amas', 207), 
        ('Booker T', 156), ('La droga', 162), ('Dakiti', 205), ('Trellas', 157), ('Sorry papi', 163), 
        ('120', 151), ('Antes que se acabe', 221), ('Cantares de Navidad', 199)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, bunny_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, bunny_id FROM current_album;


    -- Un Verano Sin TI
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('Un Verano Sin Ti', '2022-05-06', 'LP', 'https://m.media-amazon.com/images/I/8194+bM9Q1L.jpg') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Moscow Mule', 245), ('Después de la playa', 230), ('Me porto bonito', 178), ('Tití me preguntó', 243), 
        ('Un ratito', 176), ('Yo no soy celoso', 230), ('Tarot', 237), ('Neverita', 173), ('La corriente', 198),
        ('Efecto', 213), ('Party', 227), ('Aguacero', 210), ('Enséñame a bailar', 176), ('Ojitos lindos', 258), 
        ('Dos mil 16', 208), ('El apagón', 201), ('Otro atardecer', 244), ('Un coco', 196), ('Andrea', 339), 
        ('Me fui de vacaciones', 180), ('Un verano sin ti', 148), ('Agosto', 139), ('Callaíta', 250)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, bunny_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, bunny_id FROM current_album;


    -- nadie sabe lo que va a pasar mañana
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('nadie sabe lo que va a pasar mañana', '2023-10-13', 'LP', 'https://m.media-amazon.com/images/I/61TY9dEifZL.jpg') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Nadie sabe', 380), ('Monaco', 267), ('Fina', 216), ('Hibiki', 208), ('Mr. October', 190), ('Cybertruck', 192), 
        ('Vou 787', 123), ('Seda', 191), ('Gracias por nada', 176), ('Teléfono nuevo', 355), ('Baby Nueva', 241), 
        ('Mercedes Carota', 201), ('Los Pits', 251), ('Vuelve, Candy B', 266), ('Baticano', 250), ('No me quiero casar', 226), 
        ('Where She Goes', 231), ('Thunder y lightning', 218), ('Perro Negro', 163), ('Europa :(', 12), ('Acho PR', 360), 
        ('Un preview', 165)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, bunny_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, bunny_id FROM current_album;


    -- DeBÍ TiRAR MáS FOToS
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('DeBÍ TiRAR MáS FOToS', '2025-01-05', 'LP', 'https://m.media-amazon.com/images/I/81+9AGAe2oL.jpg') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('NUEVAYoL', 183), ('VOY A LLeVARTE PA PR', 156), ('BAILE INoLVIDABLE', 367), ('PERFuMITO NUEVO', 200), 
        ('WELTiTA', 187), ('VeLDÁ', 235), ('EL ClúB', 222), ('KETU TeCRÉ', 250), ('BOKeTE', 215), ('KLOuFRENS', 199), 
        ('TURiSTA', 190), ('CAFé CON RON', 228), ('PIToRRO DE COCO', 206), ('LO QUE LE PASÓ A HAWAii', 229), 
        ('EoO', 204), ('DtMF', 237), ('LA MuDANZA', 213)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, bunny_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, bunny_id FROM current_album;


    -- FACEBROOKLYN
    INSERT INTO artists (name, genre, country, bio, image_url) 
        VALUES ('FaceBrooklyn', 'Urbano', 'Chile', 'Es un cantautor y productor nacido en Viña del Mar, Chile', 'https://i.scdn.co/image/ab6761610000e5eb9bef76477b9da12aaa429fbd')
        RETURNING id INTO faceb_id;


    -- HABLANDO DE TI
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('HABLANDO DE TI', '2024-08-15', 'LP', 'https://m.media-amazon.com/images/I/A1EGTg3nD2L.jpg') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('ROMPECABEZAS', 153), ('FUGAZI', 212), ('NADIE TEKITA', 188), ('FREDDY', 192), 
        ('DEYANDEL', 187), ('FANTASSYY', 179)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, faceb_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, faceb_id FROM current_album;


    -- EASYKID
    INSERT INTO artists (name, genre, country, bio, image_url) 
        VALUES ('Easykid', 'Urbano', 'Chile', 'Joaquín Andrés Palacios Zamorano, conocido artísticamente como Easykid, es un cantante chileno de música urbana. Es conocido por su álbum conceptual Sorry, estoy en mi Darkera y su álbum más reciente I''M PART. También por canciones como Siempre pienso en ti y Vea pues.', 'https://akamai.sscdn.co/tb/letras-blog/wp-content/uploads/2025/09/3307d54-Easykid-es-parte-de-un-espacio-disruptivo-autentico-y-oscuro-1024x616.jpg')
        RETURNING id INTO easy_id;


    -- Sorry, Estoy en Mi Darkera
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('Sorry, Estoy en Mi Darkera', '2023-07-27', 'LP', 'https://m.media-amazon.com/images/I/A1GR9wDziTL.jpg') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Pa que pichea', 162), ('Siempre pienso en tí', 199), ('Coketa', 202), ('Cristina', 179), ('Tanjiro', 150), 
        ('Baby Tk', 173), ('La Tsuru', 175), ('Antibellakera', 195), ('El culto siempre gana', 139), ('Darkera', 147)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, easy_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, easy_id FROM current_album;

    -- I'M PART
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('I''M PART', '2025-06-11', 'LP', 'https://m.media-amazon.com/images/I/81fefei41JL.jpg') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Mood', 155), ('Airbag', 169), ('Whyme?', 209), ('Fentanyl', 94), ('Ya entregamos el depa', 204), 
        ('Paquepu', 130), ('Shiny', 158), ('Wtfff.', 132), ('Bruce Wayne', 137), ('Zundada de fondo', 178), 
        ('Rush', 161), ('Ojos empapados', 149), ('BANDI2', 189), ('CONTIGO/CONMIGO', 159), ('KIMYE', 193), 
        ('UN RATITO+', 109), ('LAS OLAS', 236), ('DE ELLA PARA MÍ', 168)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, easy_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, easy_id FROM current_album;

    -- +Xqa
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('+Xqa', '2021-08-19', 'LP', 'https://m.media-amazon.com/images/I/91EAsu8Wu8L.jpg') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('-10', 158), ('Wau', 203), ('Vea Pues', 182), ('Mas Na', 156), ('Lova', 130), ('Bellakera', 195), ('Silvia 120', 170)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, easy_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, easy_id FROM current_album;


    -- Visionari
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('Visionari', '2019-11-25', 'EP', 'https://m.media-amazon.com/images/I/A1POiZtcMKL.jpg') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Visionari', 215), ('Paz Koi', 193), ('Atrevia', 209), ('Tbt', 143), ('Penita', 199)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, easy_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, easy_id FROM current_album;




    INSERT INTO artists (name, genre, country, bio, image_url) 
        VALUES ('Ases Falsos', 'Rock alternativo', 'Chile', 'Ases Falsos es una banda chilena de rock alternativo formada por Cristóbal Briceño, Simón Sánchez, Martin del Real, Francisco Rojas y Daniel de la Fuente. Entre 2005 y 2011 fueron conocidos bajo el nombre de Fother Muckers.', 'https://i.scdn.co/image/ab6761610000e5eba49ade7af345b74133c1667d')
        RETURNING id INTO ases_id;

    -- Juventud Americana
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('Juventud Americana', '2012-09-03', 'LP', 'https://i.scdn.co/image/ab67616d000082c180deae578c2e2673abb839ba') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Misterios del Perú', 153), ('Salto Alto', 332), ('Pacífico', 225), ('Séptimo Cielo', 179), 
        ('Venir es fácil', 180), ('El Golfo de Adén', 179), ('No quiero que estés conmigo', 160), 
        ('Fuerza Especial', 378), ('Manantial', 303), ('Europa', 269), ('La Flor del Jazmín', 317), 
        ('Aguanieve', 200), ('Quemando', 211), ('Estudiar y Trabajar', 184), ('La Sinceridad del Cosmos', 246)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, ases_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, ases_id FROM current_album;

    -- Conducción
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('Conducción', '2014-07-01', 'LP', 'https://i.scdn.co/image/ab67616d000082c1493f8239acc8a2b7f61638c1') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Mantén la conducción', 249), ('La gran curva', 283), ('Plácidamente', 261), ('Cae la cortina', 259), 
        ('Mi ejército', 238), ('Búscate un lugar para ensayar', 246), ('Nada', 299), ('Simetría', 241), 
        ('Niña por favor', 250), ('Tora Bora', 208), ('Ivanka', 288), ('Yo no quiero volver', 355), 
        ('Al borde del cañón', 249), ('Una estrella que se mueve', 149)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, ases_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, ases_id FROM current_album;

    -- El hombre puede
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('El hombre puede', '2016-10-28', 'LP', 'https://i.scdn.co/image/ab67616d000082c1e605cac45a605953bf670e54') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Chakras', 115), ('Gehena', 158), ('Sal de ahí', 188), ('Subyugado', 249), 
        ('Más se fortalece', 329), ('Fría', 247), ('Mucho más mío', 261), ('Antes sí, ahora no', 180), 
        ('Creo que no creo', 241), ('Trato hecho', 269)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, ases_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, ases_id FROM current_album;
    


    INSERT INTO artists (name, genre, country, bio, image_url) 
    VALUES ('Coldplay', 'Pop', 'Reino Unido', 'banda británica de rock formada en Londres en 1997. Está compuesta por el vocalista y pianista Chris Martin, el guitarrista Jonny Buckland, el bajista Guy Berryman y el baterista Will Champion, además de Phil Harvey como su mánager.', 'https://i.scdn.co/image/ab6761610000e5eb1ba8fc5f5c73e7e9313cc6eb')
    RETURNING id INTO coldplay_id;

    -- Parachutes
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('Parachutes', '2000-07-10', 'LP', 'https://i.scdn.co/image/ab67616d000082c19164bafe9aaa168d93f4816a') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Don''t Panic', 137), ('Shiver', 299), ('Spies', 318), ('Sparks', 227), ('Yellow', 269), 
        ('Trouble', 270), ('Parachutes', 46), ('High Speed', 254), ('We Never Change', 249), 
        ('Everything''s Not Lost', 435)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, coldplay_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, coldplay_id FROM current_album;
  

    -- A Rush of Blood to the Head
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('A Rush of Blood to the Head', '2002-09-26', 'LP', 'https://m.media-amazon.com/images/I/71n95JlOLgL.jpg') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Politik', 318), ('In My Place', 226), ('God Put a Smile upon Your Face', 297), ('The Scientist', 309), 
        ('Clocks', 307), ('Daylight', 327), ('Green Eyes', 223), ('Warning Sign', 331), ('A Whisper', 238), 
        ('A Rush of Blood to the Head', 351), ('Ámsterdam', 319)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, coldplay_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, coldplay_id FROM current_album;

    -- X&Y
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('X&Y', '2005-06-06', 'LP', 'https://m.media-amazon.com/images/I/71SZonv3+9L.jpg') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Square One', 287), ('What If', 297), ('White Shadows', 328), ('Fix You', 294), ('Talk', 311), 
        ('X&Y', 274), ('Speed of Sound', 288), ('A Message', 285), ('Low', 332), ('The Hardest Part', 265), 
        ('Swallowed in the Sea', 238), ('Twisted Logic', 301), ('''Til Kingdom Come', 250)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, coldplay_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, coldplay_id FROM current_album;

    -- Viva la Vida or Death and All His Friends
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('Viva la Vida or Death and All His Friends', '2008-06-12', 'LP', 'https://m.media-amazon.com/images/I/9145yafeO2L.jpg') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Life in Technicolor', 149), ('Cemeteries of London', 201), ('Lost!', 235), ('42', 237), 
        ('Lovers in Japan/Reign of Love', 411), ('Yes', 426), ('Viva la vida', 241), ('Violet Hill', 222), 
        ('Strawberry Swing', 249), ('Death and All His Friends', 378)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, coldplay_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, coldplay_id FROM current_album;

    -- Mylo Xyloto
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('Mylo Xyloto', '2011-10-24', 'LP', 'https://m.media-amazon.com/images/I/91g6FY5ewrL.jpg') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Mylo Xyloto', 42), ('Hurts Like Heaven', 242), ('Paradise', 278), ('Charlie Brown', 285), 
        ('Us Against The World', 240), ('M.M.I.X.', 48), ('Every Teardrop Is a Waterfall', 241), 
        ('Major Minus', 210), ('U.F.O.', 138), ('Princess of China', 238), ('Up in Flames', 193), 
        ('A Hopeful Transmission', 33), ('Don''t Let It Break Your Heart', 234), ('Up with the Birds', 226)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, coldplay_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, coldplay_id FROM current_album;
    



    INSERT INTO artists (name, genre, country, bio, image_url) 
    VALUES ('Los Bunkers', 'Rock Pop', 'Chile', 'banda de rock chilena oriunda de Concepción, activa de 1998 a 2014, y del 2022 al presente, con una breve reunión en 2019. El grupo estuvo compuesto la mayor parte de su historia por los hermanos Francisco (voz, guitarra, teclados) y Mauricio Durán (guitarra), Álvaro (voz, guitarra) y Gonzalo López (bajo), y Mauricio Basualto (batería).', 'https://i.scdn.co/image/ab6761610000e5eb02d4a255fa006a90b4aa0605')
    RETURNING id INTO bunkers_id;

    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('Los Bunkers', '2001-04-03', 'LP', 'https://i.scdn.co/image/ab67616d000082c1d017a2bd742142174b393223') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('El detenido', 313), ('Fantasías animadas de ayer y hoy', 131), ('No sé', 304), ('Buscando cuadros', 137), 
        ('Jamás', 207), ('Yo sembré mis penas de amor en tu jardín', 168), ('Papá no llores más', 205), 
        ('Nada me importa', 214), ('Entre mis brazos', 301), ('Quiero descansar', 226)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, bunkers_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, bunkers_id FROM current_album;

    -- Canción de lejos
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('Canción de lejos', '2002-06-06', 'LP', 'https://i.scdn.co/image/ab67616d000082c1f404b4eb378a12a37fcb9d99') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Canción de lejos', 248), ('Las cosas que cambié y dejé por ti', 193), ('Sabes que...', 236), 
        ('Mañana lo voy a saber', 305), ('Canción de cerca', 248), ('Miño', 244), ('Los premios', 160), 
        ('Lo que me angustia', 223), ('Dulce final', 224), ('Pobre corazón', 239), ('Siniestra', 259), 
        ('Sacramento', 175)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, bunkers_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, bunkers_id FROM current_album;

    -- La culpa
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('La culpa', '2003-10-23', 'LP', 'https://i.scdn.co/image/ab67616d000082c18b1621e8e25425a1e70f75f4') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Canción para mañana', 195), ('No me hables de sufrir', 235), ('No necesito pensar', 224), ('Cura de espanto', 143), 
        ('Dios no sabe perder', 310), ('Culpable', 208), ('La exiliada del sur', 255), ('El día feliz', 184), 
        ('El festín de los demás', 212), ('Mariposa', 224), ('Mira lo que dicen sobre nuestro amor', 257), 
        ('Última canción', 364)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, bunkers_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, bunkers_id FROM current_album;

    -- Vida de perros
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('Vida de perros', '2003-10-23', 'LP', 'https://i.scdn.co/image/ab67616d000082c15b707596812daac339d409ad') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Ven aquí', 217), ('Nada más de mí', 193), ('Llueve sobre la ciudad', 237), ('Tú', 185), 
        ('Maribel', 246), ('Ahora que no estás', 351), ('Miéntele', 237), ('Nada es igual', 236), 
        ('Te vistes y te vas', 209), ('Vida de perros', 215), ('Hoy', 219)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, bunkers_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, bunkers_id FROM current_album;

    -- Barrio estacion
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('Barrio Estación', '2008-06-19', 'LP', 'https://i.scdn.co/image/ab67616d000082c1aebd6aea9679db9ecf83b4c6') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Coma', 279), ('Me muelen a palos', 243), ('Fiesta', 280), ('Una nube cuelga sobre mí', 238), 
        ('Andén', 257), ('Si todo esto es lo que hay', 253), ('Capablanca', 157), ('Deudas', 252), 
        ('Nada nuevo bajo el sol', 290), ('El tiempo que se va', 182), ('El mismo lugar', 182), 
        ('Tarde', 261), ('Abril', 317)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, bunkers_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, bunkers_id FROM current_album;

    -- Música libre
    WITH current_album AS (
        INSERT INTO albums (title, release_date, type, cover_url) 
        VALUES ('Música libre', '2010-10-16', 'LP', 'https://i.scdn.co/image/ab67616d000082c14a260edab0edf43df4d40844') 
        RETURNING id
    ),
    new_songs AS (
        INSERT INTO songs (title, duration) VALUES 
        ('Sueño con serpientes', 265), ('Quién fuera', 218), ('Que ya viví, que te vas', 217), 
        ('Al final de este viaje en la vida', 201), ('El necio', 208), ('Leyenda', 194), 
        ('Ángel para un final', 243), ('Santiago de Chile', 197), ('Y nada más', 133), 
        ('El día feliz que está llegando', 214), ('Pequeña serenata diurna', 129), 
        ('La era está pariendo un corazón', 286)
        RETURNING id
    ),
    relate_album AS (
        INSERT INTO tracks (album_id, song_id, track_number) 
        SELECT (SELECT id FROM current_album), id, ROW_NUMBER() OVER (ORDER BY id) FROM new_songs
    ),
    relate_artist AS (
        INSERT INTO song_artists (song_id, artist_id, role)
        SELECT id, bunkers_id, 'main' FROM new_songs
    )
    INSERT INTO album_artists (album_id, artist_id) SELECT id, bunkers_id FROM current_album;

END $$;


