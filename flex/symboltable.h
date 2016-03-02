// Author: Ersi Ni
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define ALPHABET_SIZE 26

typedef struct node {
    int data;
    struct node* link[ALPHABET_SIZE];
} node;

node* root_node = NULL;

int symboltable = 8;

node* create_node() {
    node *q = (node*) malloc(sizeof(node));
    for(int x=0;x<ALPHABET_SIZE;x++)
        q->link[x] = NULL;
    q->data = -1;
    return q;
}

int search(const char key[]) {
    if(root_node == NULL)
        root_node = create_node();
    node *q = root_node;
    int length = strlen(key);
    int level = 0;
    for(;level < length;level++) {
        int index = key[level] - 'a';
        if(q->link[index] != NULL)
            q = q->link[index];
        else
            break;
    }
    if(key[level] == '\0' && q->data != -1)
        return q->data;
    return -1;
}

void insert_node(const char key[]) {
    if (search(key) != -1) {return;}
    int length = strlen(key);
    int index;
    int level = 0;
    if(root_node == NULL)
        root_node = create_node();
    node *q = root_node;

    for(;level < length;level++) {
        index = key[level]-'a';

        if(q->link[index] == NULL) {
            q->link[index] = create_node();
        }

        q = q->link[index];
    }
    q->data = symboltable++;
}