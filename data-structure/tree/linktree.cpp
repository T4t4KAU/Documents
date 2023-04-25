#include <iostream>

enum ECCHILDSIGN {
    E_Root,       // 树根
    E_ChildLeft,  // 左孩子
    E_ChildRight  // 右孩子
};

template <typename T>
struct BinaryTreeNode {
    T data;  // 数据域
    BinaryTreeNode *leftChild;
    BinaryTreeNode *rightChild;
};

template <typename T>
class BinaryTree {
public:
    BinaryTree();
    ~BinaryTree();
public:
    // 创建节点
    BinaryTreeNode<T> *CreateNode(BinaryTreeNode<T> *parentnode,ECCHILDSIGN pointSign,const T &e);
    void ReleaseNode(BinaryTreeNode<T> *pnode);  // 释放树节点
    void CreateBTreeAccordPT(char *pstr);  // 前序遍历顺序创建BTree
public:
    // 前序遍历
    void preOrder() {
        preOrder(root);
    }
    void preOrder(BinaryTreeNode<T> *tNode) {
        if (tNode != nullptr) {
            std::cout << (char)tNode->data << " ";
            preOrder(tNode->leftChild);
            preOrder(tNode->rightChild);
        }
    }

    // 中序遍历
    void inOrder() {
        inOrder(root);
    }
    void inOrder(BinaryTreeNode<T> *tNode) {
        if (tNode != nullptr) {
            inOrder(tNode->leftChild);
            std::cout << (char)tNode->data << " ";
            inOrder(tNode->rightChild);
        }
    }

    // 后序遍历
    void postOrder() {
        postOrder(root);
    }
    void postOrder(BinaryTreeNode<T> *tNode) {
        if (tNode != nullptr) {
            postOrder(tNode->leftChild);
            postOrder(tNode->rightChild);
            std::cout << (char)tNode->data << " ";
        }
    }

private:
    BinaryTreeNode<T> *root;
    void CreateBTreeAccordPTRecu(BinaryTreeNode<T>* &tnode,char* &pstr);
};

template <class T>
BinaryTree<T>::BinaryTree() {
    root = nullptr;
}

template <class T>
BinaryTree<T>::~BinaryTree() {
    ReleaseNode(root);
}

template <class T>
void BinaryTree<T>::ReleaseNode(BinaryTreeNode<T> *pnode) {
    if (pnode != nullptr) {
        ReleaseNode(pnode->leftChild);
        ReleaseNode(pnode->rightChild);
    }
    delete pnode;
}

template <class T>
BinaryTreeNode<T> *BinaryTree<T>::CreateNode(BinaryTreeNode<T> *parentnode,ECCHILDSIGN pointSign,const T &e) {
    BinaryTreeNode<T> *tmpnode = new BinaryTreeNode<T>;
    tmpnode->data = e;
    tmpnode->leftChild = nullptr;
    tmpnode->rightChild = nullptr;

    if (pointSign == E_Root) {
        root = tmpnode;
    }
    if (pointSign == E_ChildLeft) {
        parentnode->leftChild = tmpnode;
    } else if (pointSign == E_ChildRight) {
        parentnode->rightChild = tmpnode;
    }
    return tmpnode;
}

template <class T>
void BinaryTree<T>::CreateBTreeAccordPT(char *pstr) {
    CreateBTreeAccordPTRecu(root,pstr);
}

template <class T>
void BinaryTree<T>::CreateBTreeAccordPTRecu(BinaryTreeNode<T> *&tnode,char *&pstr) {
    if (*pstr == '#') {
        tnode = nullptr;
    } else {
        tnode = new BinaryTreeNode<T>;
        tnode->data = *pstr;
        CreateBTreeAccordPTRecu(tnode->leftChild,++pstr);
        CreateBTreeAccordPTRecu(tnode->rightChild,++pstr);
    }
}

int main(void) {
    BinaryTree<int> tree;
    BinaryTreeNode<int> *rootpoint = tree.CreateNode(nullptr,E_Root,'A');
    BinaryTreeNode<int> *subpoint = tree.CreateNode(rootpoint,E_ChildLeft,'B');
    subpoint = tree.CreateNode(subpoint,E_ChildLeft,'D');
    subpoint = tree.CreateNode(rootpoint,E_ChildRight,'C');
    subpoint = tree.CreateNode(subpoint,E_ChildRight,'E');

    // tree.CreateBTreeAccordPT((char*)"ABD###C#E##"); // 直接通过前序遍历创建二叉树

    std::cout << "preorder: ";
    tree.preOrder();
    std::cout << std::endl;

    std::cout << "inorder: ";
    tree.inOrder();
    std::cout << std::endl;

    std::cout << "postorder: ";
    tree.postOrder();
    std::cout << std::endl;
    
    return 0;
}