INSERT INTO public.authors (first_name, last_name, created_at, updated_at)
VALUES ('Robert', 'Martin', now(), now()),
       ('Martin', 'Fowler', now(), now()),
       ('Erich', 'Gamma', now(), now()),
       ('Richard', 'Helm', now(), now()),
       ('Ralph', 'Johnson', now(), now()),
       ('John', 'Vlissides', now(), now()),
       ('Andrew', 'Hunt', now(), now()),
       ('David', 'Thomas', now(), now()),
       ('Brian', 'Kernighan', now(), now()),
       ( 'Dennis', 'Ritchie', now(), now());

INSERT INTO public.publishers (name, created_at, updated_at)
VALUES ('Prentice Hall', now(), now()),
       ('Addison-Wesley', now(), now()),
       ('O’Reilly Media', now(), now()),
       ('Pearson', now(), now());

INSERT INTO public.categories (name, created_at, updated_at)
VALUES ('Software Engineering', now(), now()),
       ('Programming', now(), now()),
       ('Best Practices', now(), now()),
       ('Design Patterns', now(), now()),
       ('Clean Code', now(), now());

INSERT INTO public.books (title, cover, page_count, author_id, publisher_id, publication_date, created_at,
                          updated_at, description)
VALUES ('Clean Code', '', 464, 1, 1, '2008-08-01', now(), now(),
        'Clean Code: A Handbook of Agile Software Craftsmanship by Robert C. Martin teaches developers how to write clean, maintainable, and scalable code by using good naming, simple structure, and focused functions. The book provides real-world examples of bad code transformed into good code, covering a wide range of practices that are essential for any serious software engineer.'),

       ('Clean Architecture', '', 432, 1, 2, '2017-09-20', now(), now(),
        'Clean Architecture: A Craftsman''s Guide to Software Structure and Design explores the principles behind building robust software systems. Robert C. Martin offers timeless architectural rules and explains how to separate details from policies, making systems easier to understand, maintain, and evolve. This book emphasizes the importance of dependency inversion and component separation for long-term project success.'),

       ('The Clean Coder', '', 256, 1, 2, '2011-05-13', now(), now(),
        'The Clean Coder: A Code of Conduct for Professional Programmers dives into the mindset and behaviors of a professional developer. Through anecdotes and practical advice, Robert C. Martin outlines what it means to be accountable, communicate clearly, manage time effectively, and deliver quality software. It serves as a guide for developers who want to approach coding as a disciplined craft.'),

       ('Refactoring', '', 448, 2, 2, '2018-11-19', now(), now(),
        'Refactoring: Improving the Design of Existing Code by Martin Fowler teaches how to change a software system in a way that does not alter its behavior but improves its internal structure. The book presents a catalog of common code smells and proven techniques to eliminate them. It is a vital resource for software developers and teams working to maintain code quality and agility.'),

       ('Design Patterns', '', 395, 3, 2, '1994-10-31', now(), now(),
        'Design Patterns: Elements of Reusable Object-Oriented Software, written by the "Gang of Four", is a foundational text in software engineering. It catalogs 23 classic design patterns that solve common problems in software design. Each pattern is described in detail, with UML diagrams and code examples, helping developers create flexible and reusable object-oriented designs.'),

       ('The Pragmatic Programmer', '', 352, 7, 3, '1999-10-20', now(), now(),
        'The Pragmatic Programmer: Your Journey to Mastery is a highly influential book by Andrew Hunt and David Thomas that covers the mindset and practices of successful software developers. It promotes pragmatic thinking, self-development, and the importance of communication and flexibility in coding. The book is packed with tips and best practices for writing clean, efficient, and maintainable code.'),

       ('The C Programming Language', '', 288, 9, 2, '1988-04-01', now(), now(),
        'The C Programming Language, written by Brian W. Kernighan and Dennis M. Ritchie, is the definitive guide to C, authored by one of the language\''s creators. It offers a concise and practical approach to learning the C programming language,
        including syntax, semantics, and standard libraries.It remains a must - read for those who want to understand
        system - level programming and the foundations of modern software.');

INSERT INTO public.book_category (book_id, category_id)
VALUES (1, 1),
       (1, 2),
       (1, 3),
       (2, 1),
       (2, 3),
       (3, 3),
       (4, 1),
       (4, 2),
       (5, 1),
       (5, 4),
       (6, 1),
       (6, 2),
       (7, 2);
