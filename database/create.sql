Create Table Blog {
    BlogID          int             primary key,
    UUID            uuid            not null,
    CategoryID      int             not null,        
    Title           varchar(128)    not null,
    Content         varchar(256),
    CreateAt        timestamp       not null,
    UpdatedAt       timestamp       not null,
    DeletedAt       timestamp       not null,

    Foreign key(BlogID) References Category(CategoryID)
}

Create Table BlogCategory {
    CategoryID      int             primary key,
    ProductID       int             not null,
    Name            varchar(128)    not null,
    CreateAt        timestamp       not null,
    UpdatedAt       timestamp       not null,
    DeletedAt       timestamp       not null,

    Foreign key(ProductID) References Product(CategoryID)
}


Create Table Product {
    ProductID       int             not null,
    Name            varchar(128)    not null,
    CreateAt        timestamp       not null,
    UpdatedAt       timestamp       not null,
    DeletedAt       timestamp       not null,
}
