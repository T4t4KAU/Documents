#include <iostream>

#pragma warning (disable: 4996)  
using namespace std;

//树中每个节点的定义
template <typename T> //T 代表数据元素的类型
struct BinaryTreeNode
{
	T               data;        //数据域，存放数据元素
	BinaryTreeNode* leftChild,   //左子节点指针
		* rightChild;  //右子节点指针
	int8_t        leftTag,       //左标志 = 0 表示 leftChild 指向的是左子节点，=1 表示 leftChild 指向的是前趋节点（线索）
		rightTag;      //右标志 = 0 表示 rightChild 指向的是右子节点，=1 表示 rightChild 指向的是后继节点（线索）
};

//线索二叉树的定义
template <typename T>
class ThreadBinaryTree
{
public:
	ThreadBinaryTree();  //构造函数		
	~ThreadBinaryTree(); //析构函数
public:
	//利用扩展二叉树的前序遍历序列来创建一棵二叉树
	void CreateBTreeAccordPT(char* pstr);
private:
	//利用扩展二叉树的前序遍历序列创建二叉树的递归函数
	void CreateBTreeAccordPTRecu(BinaryTreeNode<T>*& tnode, char*& pstr);//参数为引用类型，确保递归调用中对参数的改变会影响到调用者
public:
	//在二叉树中根据中序遍历序列创建线索
	void CreateThreadInBTreeAccordIO();
private:
	void CreateThreadInBTreeAccordIO(BinaryTreeNode<T>*& tnode, BinaryTreeNode<T>*& pre);//参数为引用类型

public:
	//在二叉树中根据前序遍历序列创建线索
	void CreateThreadInBTreeAccordPO();
private:
	void CreateThreadInBTreeAccordPO(BinaryTreeNode<T>*& tnode, BinaryTreeNode<T>*& pre);//参数为引用类型

public:
	//找线索二叉树（中序线索化）的第一个节点
	BinaryTreeNode<T>* GetFirst_IO()
	{
		return GetFirst_IO(root);
	}
	BinaryTreeNode<T>* GetFirst_IO(BinaryTreeNode<T>* root)
	{
		//中序遍历是“左根右顺序”
		//一直找真正的左孩子即可。
		if (root == nullptr)
			return nullptr;
		BinaryTreeNode<T>* tmpnode = root;
		while (tmpnode->leftTag == 0) //指向的是真正的左子节点
		{
			tmpnode = tmpnode->leftChild;
		}
		return tmpnode;
	}

	//找线索二叉树（中序线索化）的最后一个节点
	BinaryTreeNode<T>* GetLast_IO()
	{
		return GetLast_IO(root);
	}
	BinaryTreeNode<T>* GetLast_IO(BinaryTreeNode<T>* root)
	{
		//中序遍历是“左根右顺序”
		//一直找真正的右孩子即可。
		if (root == nullptr)
			return nullptr;
		BinaryTreeNode<T>* tmpnode = root;
		while (tmpnode->rightTag == 0) //指向的是真正的右子节点
		{
			tmpnode = tmpnode->rightChild;
		}
		return tmpnode;
	}

	//找线索二叉树（中序线索化）中当前节点的后继节点
	BinaryTreeNode<T>* GetNextPoint_IO(BinaryTreeNode<T>* currnode)
	{
		if (currnode == nullptr)
			return nullptr;
		if (currnode->rightTag == 1)//rightChild指向的是后继节点（线索）
			return currnode->rightChild;

		//如果该节点的 rightChild指向的是真正的右孩子节点，那么该怎么获得该节点的后继节点呢？
		//考虑到中序遍历的顺序为“左根右”，那么当前节点（看成根）的右子树 的第一个节点就是：
		return GetFirst_IO(currnode->rightChild); //在右子树中查找第一个要访问的节点
	}

	//找线索二叉树（中序线索化）中当前节点的前趋节点
	BinaryTreeNode<T>* GetPriorPoint_IO(BinaryTreeNode<T>* currnode)
	{
		if (currnode == nullptr)
			return nullptr;

		if (currnode->leftTag == 1)//leftChild指向的是前趋节点（线索）
			return currnode->leftChild;

		//如果该节点的 leftChild指向的是真正的左孩子节点，那么该怎么获得该节点的前趋节点呢？
		return GetLast_IO(currnode->leftChild); //在左子树中查找最后一个要访问的节点
	}


	//-------------------------------------------------------------------		
		//找线索二叉树（前序线索化）中当前节点的后继节点
	BinaryTreeNode<T>* GetNextPoint_PO(BinaryTreeNode<T>* currnode)
	{
		//根左右
		if (currnode == nullptr)
			return nullptr;

		if (currnode->rightTag == 1)
			return currnode->rightChild;

		//该节点有右孩子才能走到这里，
		//那么：如果该节点有左孩子，则该节点的后继节点必然是左孩子的第一个节点，根据根左右顺序，就是该左孩子节点
		if (currnode->leftTag == 0) //有左孩子
			return currnode->leftChild;

		//没有左孩子，而且前面已经确定了有右孩子，则根据根左右顺序
		return currnode->rightChild;
	}
	//找线索二叉树（前序线索化）中当前节点的前趋节点
	BinaryTreeNode<T>* GetPriorPoint_PO(BinaryTreeNode<T>* currnode)
	{
		if (currnode == nullptr)
			return nullptr;

		if (currnode->leftTag == 1)
			return currnode->leftChild;

		//有左孩子,但此时是无法找到当前节点的前趋节点的，除非重新进行一次前序遍历。
		// 但如果是三叉链表，也就是说可以找到当前节点的父节点，那么，就可以分几种情况
		//(1)如果当前节点没有父节点（当前节点是根节点）则根据“根左右”规则，当前节点没有前趋节点。
		//(2)如果当前节点  是  父节点的 左孩子，那么根据  根左右规则，当前节点的前趋节点就是父节点。
		//(3)如果当前节点  是  父节点的 右孩子，并且当前节点的左兄弟不存在，那么根据 根左右规则，当前节点的前趋节点就是父节点。
		//(4)如果当前节点  是  父节点的 右孩子，并且当前节点的左兄弟存在，那么根据 根左右规则，当前节点的前趋节点一定是其左兄弟这棵子树中的（按照前序遍历顺序最后一个被访问到的）节点。
		   //那么如何在“左兄弟这棵子树中”找到最后一个被访问到的节点呢？因为左兄弟节点是“左兄弟这棵子树”的根，所以：
		   //(4.1)从根节点开始，如果有右子树，则右子树中的节点肯定最后被访问。所以从根节点开始，尽可能向右子树方向找右下的节点，因为这个节点是最后被访问到的。
		   //(4.2)如果最右下角的节点还有左子树，则尽可能向左子树方向找左下的节点，这个节点是最后被访问到的。
		   //(4.3)如果这个节点又有右子树，则尽可能向右子树方向找右下的节点，因为这个节点是最后被访问到的。。。。, 这样就回到了(4.1)。
		   //(4.4)直至找到最下面一个节点，这个节点就是前序遍历中的  子树最后一个被访问到的 节点。

		//....相关代码略，你有兴趣可以自行实现
	}

	//传统中序递归遍历来遍历线索二叉树：		
	void inOrder_Org()
	{
		inOrder_Org(root);
	}
	void inOrder_Org(BinaryTreeNode<T>* tNode)  //中序遍历二叉树
	{
		if (tNode != nullptr) //若二叉树非空
		{
			//左根右顺序
			if (tNode->leftTag == 0) //是真正的左孩子
				inOrder_Org(tNode->leftChild);  //递归方式中序遍历左子树

			cout << (char)tNode->data << " "; //输出节点的数据域的值

			if (tNode->rightTag == 0) //是真正的右孩子
				inOrder_Org(tNode->rightChild); //递归方式中序遍历右子树
		}
	}

	//中序遍历按照“中序遍历序列线索化的二叉树”线索二叉树
	void inOrder_IO()
	{
		inOrder_IO(root);
	}
	void inOrder_IO(BinaryTreeNode<T>* tNode)
	{
		BinaryTreeNode<T>* tmpNode = GetFirst_IO(tNode);
		while (tmpNode != nullptr)  //从第一个节点开始一直找后继即可
		{
			cout << (char)tmpNode->data << " "; //输出节点的数据域的值
			tmpNode = GetNextPoint_IO(tmpNode);
		}
	}

	//逆向中序遍历按照“中序遍历序列线索化的二叉树”线索二叉树
	void revInOrder_IO()
	{
		revInOrder_IO(root);
	}
	void revInOrder_IO(BinaryTreeNode<T>* tNode)
	{
		BinaryTreeNode<T>* tmpNode = GetLast_IO(tNode);
		while (tmpNode != nullptr)  //从第一个节点开始一直找后继即可
		{
			cout << (char)tmpNode->data << " "; //输出节点的数据域的值
			tmpNode = GetPriorPoint_IO(tmpNode);
		}
	}

	//中序遍历序列线索化的二叉树查找某个节点(假设二叉树的节点各不相同)
	BinaryTreeNode<T>* SearchElem_IO(const T& e)
	{
		return SearchElem_IO(root, e);
	}
	BinaryTreeNode<T>* SearchElem_IO(BinaryTreeNode<T>* tNode, const T& e)
	{
		if (tNode == nullptr)
			return nullptr;
		if (tNode->data == e)  //从根开始找
			return tNode;

		//这里的代码取自于  中序遍历按照“中序遍历序列线索化的二叉树”线索二叉树inOrder_IO()的代码
		BinaryTreeNode<T>* tmpNode = GetFirst_IO(tNode);
		while (tmpNode != nullptr)  //从第一个节点开始一直找后继即可
		{
			if (tmpNode->data == e)
				return tmpNode;
			tmpNode = GetNextPoint_IO(tmpNode);
		}
		return nullptr;
	}

	//前序遍历序列线索化的二叉树查找某个节点(假设二叉树的节点各不相同)
	BinaryTreeNode<T>* SearchElem_PO(const T& e)
	{
		return SearchElem_PO(root, e);
	}
	BinaryTreeNode<T>* SearchElem_PO(BinaryTreeNode<T>* tNode, const T& e)
	{
		if (tNode == nullptr)
			return nullptr;

		BinaryTreeNode<T>* tmpNode = root; //根就是第一个节点
		while (tmpNode != nullptr)  //从第一个节点开始一直找后继即可
		{
			if (tmpNode->data == e)
				return tmpNode;
			tmpNode = GetNextPoint_PO(tmpNode);
		}
		return nullptr;
	}

	//后序遍历序列线索化的二叉树查找某个节点(假设二叉树的节点各不相同)
	BinaryTreeNode<T>* SearchElem_POSTO(const T& e)
	{
		return SearchElem_POSTO(root, e);
	}
	BinaryTreeNode<T>* SearchElem_POSTO(BinaryTreeNode<T>* tNode, const T& e)
	{
		if (tNode == nullptr)
			return nullptr;

		BinaryTreeNode<T>* tmpNode = root; //根就是最后一个节点
		while (tmpNode != nullptr)  //从最后一个节点开始一直找前趋即可
		{
			if (tmpNode->data == e)
				return tmpNode;
			tmpNode = GetPriorPoint_POSTO(tmpNode);
		}
		return nullptr;
	}

private:
	void ReleaseNode(BinaryTreeNode<T>* pnode);
private:
	BinaryTreeNode<T>* root; //树根指针	
};

//构造函数
template<class T>
ThreadBinaryTree<T>::ThreadBinaryTree()
{
	root = nullptr;
}

//析构函数
template<class T>
ThreadBinaryTree<T>::~ThreadBinaryTree()
{
	ReleaseNode(root);
};

//释放二叉树节点
template<class T>
void ThreadBinaryTree<T>::ReleaseNode(BinaryTreeNode<T>* pnode)
{
	if (pnode != nullptr)
	{
		if (pnode->leftTag == 0)
			ReleaseNode(pnode->leftChild); //只有真的需要delete的节点，才会递归调用ReleaseNode
		if (pnode->rightTag == 0)
			ReleaseNode(pnode->rightChild); //只有真的需要delete的节点，才会递归调用ReleaseNode
	}
	delete pnode;
}

//利用扩展二叉树的前序遍历序列来创建一棵二叉树
template<class T>
void ThreadBinaryTree<T>::CreateBTreeAccordPT(char* pstr)
{
	CreateBTreeAccordPTRecu(root, pstr);
}

//利用扩展二叉树的前序遍历序列创建二叉树的递归函数
template<class T>
void ThreadBinaryTree<T>::CreateBTreeAccordPTRecu(BinaryTreeNode<T>*& tnode, char*& pstr)
{
	if (*pstr == '#')
	{
		tnode = nullptr;
	}
	else
	{
		tnode = new BinaryTreeNode<T>; //创建根节点
		tnode->leftTag = tnode->rightTag = 0; //标志先给0
		tnode->data = *pstr;
		CreateBTreeAccordPTRecu(tnode->leftChild, ++pstr); //创建左子树
		CreateBTreeAccordPTRecu(tnode->rightChild, ++pstr);//创建右子树
	}
}

//在二叉树中根据中序遍历序列创建线索
template<class T>
void ThreadBinaryTree<T>::CreateThreadInBTreeAccordIO()
{
	BinaryTreeNode<T>* pre = nullptr;  //记录当前所指向的节点的前趋节点（刚开始的节点没有前趋，所以设置为nullptr）

	CreateThreadInBTreeAccordIO(root, pre);

	//注意处理最后一个节点的右孩子，因为这个右孩子还没处理
	pre->rightChild = nullptr; //这里之所以直接给nullptr，是因为中序遍历访问顺序是左根右，所以最后一个节点不可能有右孩子，否则最后一个访问的节点就会是他的右孩子。其实就算不执行这句，pre->rightChild也已经是等于nullptr的了。 
	pre->rightTag = 1; //线索化
}

template<class T>
void ThreadBinaryTree<T>::CreateThreadInBTreeAccordIO(BinaryTreeNode<T>*& tnode, BinaryTreeNode<T>*& pre)
{
	if (tnode == nullptr)
		return;

	//中序遍历序列（左根右），递归顺序非常类似于中序遍历		
	CreateThreadInBTreeAccordIO(tnode->leftChild, pre);

	if (tnode->leftChild == nullptr) //找空闲的指针域进行线索化
	{
		tnode->leftTag = 1; //线索
		tnode->leftChild = pre;  //如果leftChild ==nullptr，说明该节点没有前趋节点
	}

	//这个前趋节点的后继节点肯定是当前这个节点tnode 
	if (pre != nullptr && pre->rightChild == nullptr)
	{
		pre->rightTag = 1; //线索
		pre->rightChild = tnode;
	}
	pre = tnode; //前趋节点指针指向当前节点

	CreateThreadInBTreeAccordIO(tnode->rightChild, pre);
}

//在二叉树中根据前序遍历序列创建线索
template<class T>
void ThreadBinaryTree<T>::CreateThreadInBTreeAccordPO()
{
	BinaryTreeNode<T>* pre = nullptr;
	CreateThreadInBTreeAccordPO(root, pre);
	pre->rightChild = nullptr;
	pre->rightTag = 1;
}
template<class T>
void ThreadBinaryTree<T>::CreateThreadInBTreeAccordPO(BinaryTreeNode<T>*& tnode, BinaryTreeNode<T>*& pre)
{
	if (tnode == nullptr)
		return;

	//前遍历序列（根左右），递归顺序非常类似于前序遍历
	if (tnode->leftChild == nullptr) //找空闲的指针域进行线索化
	{
		tnode->leftTag = 1; //线索
		tnode->leftChild = pre;  //如果leftChild ==nullptr，说明该节点没有前趋节点
	}

	//这个前趋节点的后继节点肯定是当前这个节点tnode 
	if (pre != nullptr && pre->rightChild == nullptr)
	{
		pre->rightTag = 1; //线索
		pre->rightChild = tnode;
	}
	pre = tnode; //前趋节点指针指向当前节点

	if (tnode->leftTag == 0) //当leftChild是真正的子节点而不是线索化后的前趋节点时
		CreateThreadInBTreeAccordPO(tnode->leftChild, pre);

	if (tnode->rightTag == 0) //当rightChild是真正的子节点而不是线索化后的后继趋节点时
		CreateThreadInBTreeAccordPO(tnode->rightChild, pre);
}


int main()
{

	ThreadBinaryTree<int> mythreadtree;
	mythreadtree.CreateBTreeAccordPT((char*)"ABD#G##EH###C#F##");  //利用扩展二叉树的前序遍历序列创建二叉树
	//对二叉树进行线索化(根据中序遍历序列创建线索)
	mythreadtree.CreateThreadInBTreeAccordIO();
	//mythreadtree.CreateThreadInBTreeAccordPO();

	//--------------------
	mythreadtree.inOrder_Org(); //传统中序递归遍历
	cout << endl;
	mythreadtree.inOrder_IO();
	cout << endl;
	mythreadtree.revInOrder_IO();
	cout << endl;

	//--------------------
	int val = 'B';
	BinaryTreeNode<int>* p = mythreadtree.SearchElem_IO(val);
	if (p != nullptr)
	{
		cout << "找到了值为" << (char)val << "的节点" << endl;

		//顺便找下后继和前趋节点
		BinaryTreeNode<int>* nx = mythreadtree.GetNextPoint_IO(p);
		if (nx != nullptr)
			cout << "后继节点值为" << (char)nx->data << "." << endl;

		BinaryTreeNode<int>* pr = mythreadtree.GetPriorPoint_IO(p);
		if (pr != nullptr)
			cout << "前趋节点值为" << (char)pr->data << "." << endl;

	}
	else
		cout << "没找到值为" << (char)val << "的节点" << endl;

	return 0;
}